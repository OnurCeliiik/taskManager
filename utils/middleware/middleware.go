package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"task-manager/internal/task"

	"task-manager/utils/error"
	"task-manager/utils/jwt"
	"task-manager/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.RWMutex
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization header is missing",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid authorization format",
			})
			return
		}

		tokenString := parts[1]

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "invalid or expired token",
			})
			return
		}

		c.Set("userID", claims.UserID.String())
		c.Set("role", claims.Role)

		c.Next()

	}
}

func TaskOwnershipMiddleware(taskService task.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, error.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid user ID format",
			})
			return
		}

		taskIDStr := c.Param("id")
		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, error.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid task ID format",
			})
			return
		}

		taskResp, err := taskService.GetTaskByID(taskID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, error.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "task not found",
			})
			return
		}

		if taskResp.UserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, error.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "you do not have permission to access this task",
			})
			return
		}

		c.Next()
	}
}

// ensures only admins can access certain endpoints
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetString("role")
		if roleStr != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, error.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "admin access required",
			})
			return
		}
		c.Next()
	}
}

// restricts routes to a specific role (or admin)
func AuthorizationMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		if userIDStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "user not authenticated",
			})
			return
		}

		roleStr := c.GetString("role")
		if roleStr != requiredRole && roleStr != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, error.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "insufficient permissions for this action",
			})
			return
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.RequestURI
		clientIP := c.ClientIP()

		logger.Info("HTTP request completed",
			"method", method,
			"path", path,
			"status_code", statusCode,
			"duration_ms", duration.Milliseconds(),
			"client_ip", clientIP,
		)
	}
}

func getRateLimiter(ip string) *rate.Limiter {
	mu.RLock()
	limiter, exists := limiters[ip]
	mu.RUnlock()

	if !exists {
		// 10 requests per second, burst of 20
		limiter = rate.NewLimiter(10, 20)
		mu.Lock()
		limiters[ip] = limiter
		mu.Unlock()
	}

	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := getRateLimiter(clientIP)

		if !limiter.Allow() {
			logger.Warn("rate limit exceeded", "ip", clientIP, "path", c.Request.RequestURI)
			c.JSON(http.StatusTooManyRequests, error.ErrorResponse{
				Code:    http.StatusTooManyRequests,
				Message: "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func StrictRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mu.RLock()
		limiter, exists := limiters[clientIP+"_auth"]
		mu.RUnlock()

		if !exists {
			limiter = rate.NewLimiter(5, 10)
			mu.Lock()
			limiters[clientIP+"_auth"] = limiter
			mu.Unlock()
		}

		if !limiter.Allow() {
			logger.Warn("auth rate limit exceeded", "ip", clientIP)
			c.JSON(http.StatusTooManyRequests, error.ErrorResponse{
				Code:    http.StatusTooManyRequests,
				Message: "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
