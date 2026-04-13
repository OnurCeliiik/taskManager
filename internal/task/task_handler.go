package task

import (
	"net/http"
	"strconv"

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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

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

func (h *TaskHandler) ListTasks(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	status := c.Query("status")
	category := c.Query("category")

	tasks, err := h.service.ListTasks(userID, limit, offset, status, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list tasks",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid task ID",
		})
		return
	}

	task, err := h.service.GetTaskByID(taskID)
	if err != nil {
		if err == ErrTaskNotFound {
			c.JSON(http.StatusNotFound, error.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid task ID",
		})
		return
	}

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	task, err := h.service.UpdateTask(taskID, userID, req)
	if err != nil {
		if err == ErrTaskNotFound {
			c.JSON(http.StatusNotFound, error.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Task not found",
			})
			return
		}
		if err == ErrInvalidTaskStatus {
			c.JSON(http.StatusBadRequest, error.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid status value",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update task",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid task ID",
		})
		return
	}

	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, error.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}

	if err := h.service.DeleteTask(taskID, userID); err != nil {
		if err == ErrTaskNotFound {
			c.JSON(http.StatusNotFound, error.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, error.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete task",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
