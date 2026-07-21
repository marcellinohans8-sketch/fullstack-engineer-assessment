package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"backend/config"
	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type taskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	Assignee    string `json:"assignee"`
}

type taskFilters struct {
	Status   string
	Keyword  string
	Assignee string
	Page     int
	Limit    int
	Sort     string
}

var allowedTaskSorts = map[string]string{
	"created_at asc":  "created_at asc",
	"created_at desc": "created_at desc",
	"createdAt asc":   "created_at asc",
	"createdAt desc":  "created_at desc",
	"updated_at asc":  "updated_at asc",
	"updated_at desc": "updated_at desc",
	"updatedAt asc":   "updated_at asc",
	"updatedAt desc":  "updated_at desc",
	"title asc":       "title asc",
	"title desc":      "title desc",
	"status asc":      "status asc",
	"status desc":     "status desc",
	"assignee asc":    "assignee asc",
	"assignee desc":   "assignee desc",
}

func CreateTask(c *gin.Context) {
	var input taskRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	task := models.Task{
		Title:       strings.TrimSpace(input.Title),
		Description: input.Description,
		Status:      normalizeStatus(input.Status),
		Assignee:    input.Assignee,
	}

	if err := config.DB.Create(&task).Error; err != nil {
		handleTaskWriteError(c, err, "Failed to create task")
		return
	}

	_ = helpers.DeleteCache("tasks:*")

	helpers.Success(c, http.StatusCreated, "Task created successfully", task)
}

func GetTasks(c *gin.Context) {
	var tasks []models.Task

	filters, err := parseTaskFilters(c)
	if err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	cacheKey := buildTasksCacheKey(filters)
	var cached map[string]interface{}

	if err := helpers.GetCache(cacheKey, &cached); err == nil {
		c.JSON(http.StatusOK, cached)
		return
	}

	offset := (filters.Page - 1) * filters.Limit
	query := config.DB.Model(&models.Task{})

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.Keyword != "" {
		keywordPattern := "%" + filters.Keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", keywordPattern, keywordPattern)
	}

	if filters.Assignee != "" {
		query = query.Where("assignee = ?", filters.Assignee)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Failed to count tasks", err.Error())
		return
	}

	if err := query.
		Order(filters.Sort).
		Limit(filters.Limit).
		Offset(offset).
		Find(&tasks).Error; err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Failed to get tasks", err.Error())
		return
	}

	response := gin.H{
		"success": true,
		"data":    tasks,
		"pagination": gin.H{
			"page":  filters.Page,
			"limit": filters.Limit,
			"total": total,
		},
	}

	_ = helpers.SetCache(cacheKey, response)

	c.JSON(http.StatusOK, response)
}

func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		helpers.Error(c, http.StatusBadRequest, "Invalid task id", "id must be a positive integer")
		return
	}

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.Error(c, http.StatusNotFound, "Task not found", "task does not exist")
			return
		}

		helpers.Error(c, http.StatusInternalServerError, "Failed to find task", err.Error())
		return
	}

	var input taskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	task.Title = strings.TrimSpace(input.Title)
	task.Description = input.Description
	task.Status = normalizeStatus(input.Status)
	task.Assignee = input.Assignee

	if err := config.DB.Save(&task).Error; err != nil {
		handleTaskWriteError(c, err, "Failed to update task")
		return
	}

	_ = helpers.DeleteCache("tasks:*")

	helpers.Success(c, http.StatusOK, "Task updated successfully", task)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		helpers.Error(c, http.StatusBadRequest, "Invalid task id", "id must be a positive integer")
		return
	}

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.Error(c, http.StatusNotFound, "Task not found", "task does not exist")
			return
		}

		helpers.Error(c, http.StatusInternalServerError, "Failed to find task", err.Error())
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Failed to delete task", err.Error())
		return
	}

	_ = helpers.DeleteCache("tasks:*")

	helpers.Success(c, http.StatusOK, "Task deleted successfully", nil)
}

func parseTaskFilters(c *gin.Context) (taskFilters, error) {
	page, err := parsePositiveInt(c.DefaultQuery("page", "1"), "page")
	if err != nil {
		return taskFilters{}, err
	}

	limit, err := parsePositiveInt(c.DefaultQuery("limit", "10"), "limit")
	if err != nil {
		return taskFilters{}, err
	}

	if limit > 100 {
		limit = 100
	}

	status := strings.TrimSpace(c.Query("status"))
	if status != "" && !isValidStatus(status) {
		return taskFilters{}, fmt.Errorf("status must be one of todo, in_progress, done")
	}

	sortParam := strings.TrimSpace(c.DefaultQuery("sort", "created_at desc"))
	sort, ok := allowedTaskSorts[sortParam]
	if !ok {
		return taskFilters{}, fmt.Errorf("sort must be one of created_at, updated_at, title, status, assignee with asc or desc")
	}

	return taskFilters{
		Status:   status,
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		Assignee: strings.TrimSpace(c.Query("assignee")),
		Page:     page,
		Limit:    limit,
		Sort:     sort,
	}, nil
}

func parsePositiveInt(value string, field string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", field)
	}

	return parsed, nil
}

func buildTasksCacheKey(filters taskFilters) string {
	values := url.Values{}
	values.Set("assignee", filters.Assignee)
	values.Set("keyword", filters.Keyword)
	values.Set("limit", strconv.Itoa(filters.Limit))
	values.Set("page", strconv.Itoa(filters.Page))
	values.Set("sort", filters.Sort)
	values.Set("status", filters.Status)

	return "tasks:" + values.Encode()
}

func normalizeStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return "todo"
	}

	return status
}

func isValidStatus(status string) bool {
	return status == "todo" || status == "in_progress" || status == "done"
}

func handleTaskWriteError(c *gin.Context, err error, fallbackMessage string) {
	var mysqlErr *mysqlDriver.MySQLError

	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		helpers.Error(c, http.StatusConflict, "Task title already exists", "title must be unique")
		return
	}

	errText := strings.ToLower(err.Error())
	if strings.Contains(errText, "unique") || strings.Contains(errText, "duplicate") {
		helpers.Error(c, http.StatusConflict, "Task title already exists", "title must be unique")
		return
	}

	helpers.Error(c, http.StatusInternalServerError, fallbackMessage, err.Error())
}
