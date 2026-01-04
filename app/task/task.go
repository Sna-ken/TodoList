package task

import (
	"context"
	"strconv"

	"github.com/Sna-ken/hellogo/config"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func CreateTask(ctx context.Context, c *app.RequestContext) {
	userid, exsist := c.Get("user_id")
	if !exsist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}

	userID := userid.(uint)

	var req CreateTaskReq
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: consts.StatusBadRequest,
			Msg:    "invalid input",
		})
		return
	}

	if req.Title == "" {
		c.JSON(consts.StatusBadRequest, Response{
			Status: consts.StatusBadRequest,
			Msg:    "expected title",
		})
		return
	}

	task := Task{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		DueAt:   req.DueAt,
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "failed to create task",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data:   task,
	})
}

func updateTaskCompleted(userID uint, taskID *uint, completed bool) (int64, error) {
	db := config.DB.Model(&Task{}).Where("user_id = ?", userID)

	if taskID != nil {
		db = db.Where("id = ?", *taskID)
	}

	rsl := db.Update("completed", completed)
	return rsl.RowsAffected, rsl.Error
}

func CompleteSingleTask(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: consts.StatusBadRequest,
			Msg:    "invalid task id",
		})
		return
	}
	taskID := uint(id64)

	affected, err := updateTaskCompleted(userID, &taskID, true)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "update failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]any{
			"updated": affected,
			"status":  "completed",
		},
	})
}

func CompleteAllTasks(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	affected, err := updateTaskCompleted(userID, nil, true)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "update failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]any{
			"updated": affected,
			"status":  "completed",
		},
	})
}

func UncompleteSingleTask(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: consts.StatusBadRequest,
			Msg:    "invalid task id",
		})
		return
	}
	taskID := uint(id64)

	affected, err := updateTaskCompleted(userID, &taskID, false)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "update failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]any{
			"updated": affected,
			"status":  "uncompleted",
		},
	})
}

func UncompleteAllTasks(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	affected, err := updateTaskCompleted(userID, nil, false)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "update failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]any{
			"updated": affected,
			"status":  "uncompleted",
		},
	})
}

func ListTasks(ctx context.Context, c *app.RequestContext) {
	userid, exsist := c.Get("user_id")
	if !exsist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	completedStr := c.Query("completed")
	keyword := c.Query("keyword")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	db := config.DB.Where("user_id = ?", userID)

	if completedStr != "" {
		completed, err := strconv.ParseBool(completedStr)
		if err != nil {
			c.JSON(consts.StatusBadRequest, Response{
				Status: consts.StatusBadRequest,
				Msg:    "invalid completed value",
			})
			return
		}
		db = db.Where("completed = ?", completed)
	}

	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", like, like)
	}

	var total int64
	if err := db.Model(&Task{}).Count(&total).Error; err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "count failed",
		})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size
	db = db.Offset(offset).Limit(size)

	var tasks []Task

	if err := db.Order("created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "query failed",
		})
		return
	}

	resp := Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: TaskListData{
			Items: tasks,
			Total: total,
		},
	}

	c.JSON(consts.StatusOK, resp)
}

func deleteTasks(userID uint, taskID *uint, completed *bool) (int64, error) {
	db := config.DB.Where("user_id = ?", userID)
	if taskID != nil {
		db = db.Where("id = ?", *taskID)
	}
	if completed != nil {
		db = db.Where("completed = ?", *completed)
	}

	result := db.Delete(&Task{})
	return result.RowsAffected, result.Error
}

func DeleteSingleTask(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, Response{
			Status: consts.StatusBadRequest,
			Msg:    "invalid task id",
		})
		return
	}
	taskID := uint(id64)

	affected, err := deleteTasks(userID, &taskID, nil)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "delete failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]int64{
			"deleted": affected,
		},
	})
}

func DeleteCompleteAllTasks(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	completed := true
	affected, err := deleteTasks(userID, nil, &completed)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "delete failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]int64{
			"deleted": affected,
		},
	})
}

func DeleteUncompleteAllTasks(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "UserID not found",
		})
		return
	}
	userID := userid.(uint)

	completed := false
	affected, err := deleteTasks(userID, nil, &completed)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "delete failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]int64{
			"deleted": affected,
		},
	})
}

func DeleteAllTasks(ctx context.Context, c *app.RequestContext) {
	userid, exist := c.Get("user_id")
	if !exist {
		c.JSON(consts.StatusUnauthorized, Response{
			Status: consts.StatusUnauthorized,
			Msg:    "User ID not found",
		})
		return
	}
	userID := userid.(uint)

	affected, err := deleteTasks(userID, nil, nil)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			Status: consts.StatusInternalServerError,
			Msg:    "delete failed",
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data: map[string]int64{
			"deleted": affected,
		},
	})
}
