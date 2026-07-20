package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"backend/config"
	"backend/helpers"
	"backend/models"

	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := config.DB.Create(&task).Error; err != nil {
		var mysqlErr *mysqlDriver.MySQLError

		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			helpers.Error(c, http.StatusConflict, "Task title already exists", err.Error())
			return
		}

		helpers.Error(c, http.StatusInternalServerError, "Failed to create task", err.Error())
		return
	}

	_ = helpers.DeleteCache("tasks:*")

	helpers.Success(c, http.StatusCreated, "Task created successfully", task)
}

func GetTasks(c *gin.Context) {
	var tasks []models.Task

	status := c.Query("status")
	keyword := c.Query("keyword")
	assignee := c.Query("assignee")

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sort := c.DefaultQuery("sort", "created_at desc")

	cacheKey := fmt.Sprintf(
		"tasks:%s:%s:%s:%s:%s:%s",
		status,
		keyword,
		assignee,
		page,
		limit,
		sort,
	)

	var cached map[string]interface{}

	if err := helpers.GetCache(cacheKey, &cached); err == nil {
		c.JSON(http.StatusOK, cached)
		return
	}

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	if pageInt < 1 {
		pageInt = 1
	}

	if limitInt < 1 {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt

	query := config.DB.Model(&models.Task{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if keyword != "" {
		query = query.Where(
			"title LIKE ? OR description LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	if assignee != "" {
		query = query.Where("assignee = ?", assignee)
	}

	var total int64
	query.Count(&total)

	fmt.Println("📦 Fetching from MySQL")

	if err := query.
		Order(sort).
		Limit(limitInt).
		Offset(offset).
		Find(&tasks).Error; err != nil {

		helpers.Error(c, http.StatusInternalServerError, "Failed to get tasks", err.Error())
		return
	}

	response := gin.H{
		"success": true,
		"data": tasks,
		"pagination": gin.H{
			"page":  pageInt,
			"limit": limitInt,
			"total": total,
		},
	}

	_ = helpers.SetCache(cacheKey, response)

	c.JSON(http.StatusOK, response)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		helpers.Error(c, http.StatusNotFound, "Task not found", err.Error())
		return
	}

	var input models.Task

	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.Assignee = input.Assignee

	if err := config.DB.Save(&task).Error; err != nil {

		var mysqlErr *mysqlDriver.MySQLError

		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			helpers.Error(c, http.StatusConflict, "Task title already exists", err.Error())
			return
		}

		helpers.Error(c, http.StatusInternalServerError, "Failed to update task", err.Error())
		return
	}

	_ = helpers.DeleteCache("tasks:*")

	helpers.Success(c, http.StatusOK, "Task updated successfully", task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		helpers.Error(c, http.StatusNotFound, "Task not found", err.Error())
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		helpers.Error(c, http.StatusInternalServerError, "Failed to delete task", err.Error())
		return
	}

	_ = helpers.DeleteCache("tasks:*")

	helpers.Success(c, http.StatusOK, "Task deleted successfully", nil)
}