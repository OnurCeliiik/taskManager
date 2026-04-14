package routes

import (
	"task-manager/internal/auth"
	"task-manager/internal/health"
	"task-manager/internal/task"
	"task-manager/internal/users"
	"task-manager/utils/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	// Initialize user repository, service, and handler
	userRepository := users.NewUserRepository(db)
	userService := users.NewUserService(userRepository)
	userHandler := users.NewUserHandler(userService)

	// Initialize auth service and handler
	authService := auth.NewAuthService(userRepository)
	authHandler := auth.NewAuthHandler(authService)

	// Initialize task repository, service, and handler
	taskRepository := task.NewTaskRepository(db)
	taskService := task.NewTaskService(taskRepository)
	taskHandler := task.NewTaskHandler(taskService)

	// HealthCheck
	healthHandler := health.NewHealthHandler(db)

	// This creates a new Gin router with default middleware, although I don't know how to implement middleware.
	r := gin.Default()

	// the health check, can be improved further down the line
	r.GET("/healthz", healthHandler.Check)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/login", authHandler.LoginUser)
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.PUT("/:id", userHandler.UpdateUser)
			users.GET("/:id", userHandler.GetUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware())
		{
			tasks.POST("/create", taskHandler.CreateTask)
			tasks.GET("/", taskHandler.ListTasks)

			tasks.PUT("/:id", middleware.TaskOwnershipMiddleware(taskService), taskHandler.UpdateTask)
			tasks.GET("/:id", middleware.TaskOwnershipMiddleware(taskService), taskHandler.GetTask)
			tasks.DELETE("/:id", middleware.TaskOwnershipMiddleware(taskService), taskHandler.DeleteTask)

		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			// Future admin endpoints (user management, analytics, etc...)
			// admin.GET("/users", adminHandler.ListAllUsers)
			// admin.DELETE("/users/:id", adminHandler.DeleteUserAsAdmin)
		}
	}

	return r
}
