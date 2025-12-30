package user

import (
	"context"
	"time"

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
		c.JSON(consts.StatusBadRequest, "Invalid input") //无效输入
		return
	}

	hashedPassword, err := utils.HandPassword(user.Password)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, "Error hashing password") //加密错误
		return
	}

	user.Password = hashedPassword

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(consts.StatusInternalServerError, "Error creating user") //创建失败
		return
	}

	c.JSON(consts.StatusOK, "User registered successfully")
}

func Login(ctx context.Context, c *app.RequestContext) {
	//凭证
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&credentials); err != nil {
		c.JSON(consts.StatusBadRequest, "Invalid input")
		return
	}

	var user User

	if err := config.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(consts.StatusUnauthorized, "Invalid credentials") //未找到用户名
			return
		}

		c.JSON(consts.StatusInternalServerError, "Error querying user") //查询失败
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(consts.StatusUnauthorized, "Invalid credentials") //密码不一致
		return
	}

	token, err := jwt.GenerateJWT(user.ID) //生成JWT
	if err != nil {
		c.JSON(consts.StatusInternalServerError, "Error generating token") //token生成失败
		return
	}

	c.JSON(consts.StatusOK, map[string]string{"token": token})
}

func Profile(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, "User ID not found")
		return
	}

	var user User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(consts.StatusNotFound, "User not found")
			return
		}
		c.JSON(consts.StatusInternalServerError, "Erroe fetching user")
		return
	}
	profile := UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
	//这样就不会传回密码了
	c.JSON(consts.StatusOK, profile)
}
