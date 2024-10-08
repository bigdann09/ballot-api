package controllers

import (
	"net/http"

	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"
	"github.com/gin-gonic/gin"
)

func GetAllTasksController(c *gin.Context) {
	user, err := utils.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
	}

	tasks, err := models.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if tasks == nil {
		tasks = []*utils.TaskAPI{}
	} else {
		for _, task := range tasks {
			task.Completed = models.CheckTask(int64(user.TGID), task.ID)
		}
	}

	// update last activity
	models.UpdateLastActivity(user.ID)

	c.JSON(http.StatusOK, tasks)
}

func StoreTaskController(c *gin.Context) {
	var task utils.TaskCreateApiRequest
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
		return
	}

	// validate task data
	validated := utils.ValidateData(task, []string{"ID", "Link", "Completed"})
	if len(validated) > 0 {
		c.JSON(http.StatusBadRequest, validated)
		return
	}

	if found := models.CheckTaskByName(task.Name); found {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Task already exists",
		})
		return
	}

	err := models.NewTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task created successfully",
	})
}

func MarkTaskCompleteController(c *gin.Context) {
	uuid := c.Param("uuid")

	user, err := utils.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
	}

	task, err := models.GetTaskByUUID(uuid, user.TGID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if task.Completed {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Task already completed",
		})
		return
	}

	err = models.CompleteTask(user.TGID, task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// update last activity
	models.UpdateLastActivity(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Task completed.",
	})
}
