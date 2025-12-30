package config

var Mysql = struct {
	USERNAME string
	PASSWORD string
	HOST     string
	PORT     string
	NAME     string
}{
	USERNAME: "root",
	PASSWORD: "root",
	HOST:     "127.0.0.1",
	PORT:     "3306",
	NAME:     "todolist",
}
