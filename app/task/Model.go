package task

import "time"

type Task struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"-" gorm:"index;not null"`

	Title   string `json:"title" gorm:"type:varchar(255);not null"`
	Content string `json:"content" gorm:"type:text"`

	Completed bool `json:"completed" gorm:"default:false"`

	CreatedAt time.Time  `json:"created_at"`
	DueAt     *time.Time `json:"due_at"` //指针可以为null
}

type CreateTaskReq struct {
	Title   string     `json:"title"`
	Content string     `json:"content"`
	DueAt   *time.Time `json:"due_at"`
}

type TaskListData struct {
	Items []Task `json:"items"`
	Total int64  `json:"total"`
}

type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
