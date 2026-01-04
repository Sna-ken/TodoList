package user

import (
	"context"
	"time"

	"github.com/Sna-ken/hellogo/app/task"
	"github.com/Sna-ken/hellogo/config"
	"github.com/Sna-ken/hellogo/pkg/jwt"
	"github.com/Sna-ken/hellogo/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserProfile struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	var user User
	if err := c.Bind(&user); err != nil {
		c.JSON(consts.StatusBadRequest, task.Response{
			Status: consts.StatusBadRequest,
			Msg:    "Invalid input",
		})
		return
	}

	hashedPassword, err := utils.HandPassword(user.Password)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, task.Response{
			Status: consts.StatusInternalServerError,
			Msg:    "Error hashing password",
		})
		return
	}

	user.Password = hashedPassword

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(consts.StatusInternalServerError, task.Response{
			Status: consts.StatusInternalServerError,
			Msg:    "Error creating user",
		})
		return
	}

	c.JSON(consts.StatusOK, task.Response{
		Status: consts.StatusOK,
		Msg:    "User registered successfully",
	})
}

func LoginLogin(ctx context.Context, c *app.RequestContext) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&credentials); err != nil {
		c.JSON(consts.StatusBadRequest, task.Response{
			Status: consts.StatusBadRequest,
			Msg:    "Invalid input",
		})
		return
	}

	var user User

	if err := config.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(consts.StatusUnauthorized, task.Response{
				Status: consts.StatusUnauthorized,
				Msg:    "Invalid credentials",
			})
			return
		}

		c.JSON(consts.StatusInternalServerError, task.Response{
			Status: consts.StatusInternalServerError,
			Msg:    "Error querying user",
		})
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(consts.StatusUnauthorized, task.Response{
			Status: consts.StatusUnauthorized,
			Msg:    "Invalid credentials",
		})
		return
	}

	token, err := jwt.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, task.Response{
			Status: consts.StatusInternalServerError,
			Msg:    "Error generating token",
		})
		return
	}

	c.JSON(consts.StatusOK, task.Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data:   map[string]string{"token": token},
	})
}

func Profile(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, task.Response{
			Status: consts.StatusUnauthorized,
			Msg:    "User ID not found",
		})
		return
	}

	var user User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(consts.StatusNotFound, task.Response{
				Status: consts.StatusNotFound,
				Msg:    "User not found",
			})
			return
		}
		c.JSON(consts.StatusInternalServerError, task.Response{
			Status: consts.StatusInternalServerError,
			Msg:    "Error fetching user",
		})
		return
	}
	profile := UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(consts.StatusOK, task.Response{
		Status: consts.StatusOK,
		Msg:    "ok",
		Data:   profile,
	})
}
