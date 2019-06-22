package main

import (
	"github.com/oniikal3/school/todo"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	r.GET("/api/todos/:id", todo.GetHandler)
	r.POST("/api/todos", todo.PostHandler)
	r.DELETE("/api/todos/:id", todo.DeleteHandler)
	r.Run(":1234")
}
