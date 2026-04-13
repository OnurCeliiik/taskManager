package routes

import (
	"task-manager/internal/auth"
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

	taskRepository := task.NewTaskRepository(db)
	taskService := task.NewTaskService(taskRepository)
	taskHandler := task.NewTaskHandler(taskService)

	// This creates a new Gin router with default middleware, although I don't know how to implement middleware.
	r := gin.Default()

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
			// users.PUT("/:id", userHandler.UpdateUser)
			users.GET("/:id", userHandler.GetUser)
			// users.DELETE("/:id", userHandler.DeleteUser)
		}

		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware())
		{
			tasks.POST("/create", taskHandler.CreateTask)
			/*	tasks.GET("/", taskHandler.ListTasks)

				tasks.PUT("/:id", middleware.TaskOwnershipMiddleware(taskService), taskHandler.UpdateTask)
				tasks.GET("/:id", middleware.TaskOwnershipMiddleware(taskService), taskHandler.GetTask)
				tasks.DELETE("/:id", middleware.TaskOwnershipMiddleware(taskService), taskHandler.DeleteTask)
			*/
		}
	}

	return r
}
