package main

import (
	"github.com/Sna-ken/hellogo/app/task"
	"github.com/Sna-ken/hellogo/app/user"
	"github.com/Sna-ken/hellogo/config"
	"github.com/Sna-ken/hellogo/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	config.InitDB()

	if config.DB == nil {
		panic("Database connection failed")
	}

	h := server.Default()

	config.DB.AutoMigrate(&user.User{})
	config.DB.AutoMigrate(&task.Task{})

	h.POST("/register", user.Register)
	h.POST("/login", user.Login)

	group := h.Group("/protected") //受保护路由,需要toeken登录
	group.Use(middleware.JWTAuth())
	{
		group.GET("/profile", user.Profile)
		//增
		group.POST("/tasks", task.CreateTask)
		//改
		group.PUT("/tasks/:id/complete", task.CompleteSingleTask)
		group.PUT("/tasks/:id/uncomplete", task.UncompleteSingleTask)
		group.PUT("/tasks/complete", task.CompleteAllTasks)
		group.PUT("/tasks/uncomplete", task.UncompleteAllTasks)
		//查
		group.GET("/tasks", task.ListTasks)
		//删
		group.DELETE("/tasks/:id", task.DeleteSingleTask)
		group.DELETE("/tasks", task.DeleteAllTasks)
		group.DELETE("/tasks/completed", task.DeleteCompleteAllTasks)
		group.DELETE("/tasks/uncompleted", task.DeleteUncompleteAllTasks)
	}

	h.Spin()
}
