package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	r.GET("/students", getStudentHandler)
	r.POST("/students", insertStudentHandler)
	r.GET("/api/todos", getTodos)
	r.Run(":1234")
}

type Student struct {
	Name string `json:"name"`
	ID   int    `json:"student_id"`
}

var students = map[int]Student{
	1: Student{Name: "Champ", ID: 1},
	2: Student{Name: "Heroes", ID: 2},
}

func getStudentHandler(c *gin.Context) {
	ss := []Student{}
	for _, v := range students {
		ss = append(ss, v)
	}
	c.JSON(http.StatusOK, ss)
}

func insertStudentHandler(c *gin.Context) {

	// Using Original Golang
	// body := c.Request.Body
	// b, err := ioutil.ReadAll(body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stu := Student{}
	// jerr := json.Unmarshal(b, &stu)
	// if err != nil {
	// 	log.Fatal(jerr)
	// }

	// Using Gin method
	s := Student{}
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id := len(students)
	id++
	s.ID = id
	students[id] = s

	response := gin.H{
		"students": fmt.Sprintf("New student: %v", s),
	}
	c.JSON(http.StatusOK, response)
}

func getTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	stmt, _ := db.Prepare("SELECT id, title, status FROM todos ORDER by id")
	rows, _ := stmt.Query()

	todos := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
		}
		todos = append(todos, t)
	}

	c.JSON(http.StatusOK, todos)
}

type Todo struct {
	ID     int
	Title  string
	Status string
}
