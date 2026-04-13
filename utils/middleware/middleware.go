package middleware

import (
	"net/http"
	"strings"

	"task-manager/internal/task"

	"task-manager/utils/error"
	"task-manager/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		c.Set("userID", claims.UserID)

		c.Next()

	}
}

func TaskOwnershipMiddleware(taskService task.TaskService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement logic to check if the user owns the task
		// You can extract the user ID from the context and compare it with the task's owner ID
		// If the user does not own the task, return an unauthorized error

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

		task, err := task.TaskService.GetTaskByID(taskService, taskID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, error.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "task not found",
			})
			return
		}

		if task.UserID != userID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, error.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "you do not have permission to access this task",
			})
			return
		}

		c.Next()
	}
}
