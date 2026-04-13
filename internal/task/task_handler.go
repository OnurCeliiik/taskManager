package task

import (
	"net/http"

	"task-manager/utils/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest

	// 1. Bind the request body to the struct and validate it
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	// 2. Extract user ID from context
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	//
	task, err := h.service.CreateTask(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}
