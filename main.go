package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	r.GET("/students", getStudentHandler)
	r.POST("/students", insertStudentHandler)
	r.GET("/api/todos", getTodos)
	r.GET("/api/todos/:id", getTodoByIdHandler)
	r.POST("/api/todos", postTodoHandler)
	r.DELETE("/api/todos/:id", deleteTodoByIdHandler)
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

type Todo struct {
	ID     int
	Title  string
	Status string
}

func getTodos(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	stmt, _ := db.Prepare("SELECT id, title, status FROM todos ORDER by id;")
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

func dbConnect() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	return db
}

func getTodoByIdHandler(c *gin.Context) {
	db := dbConnect()
	id, _ := strconv.Atoi(c.Param("id"))
	stmt, _ := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1;")
	row := stmt.QueryRow(id)
	t := Todo{}
	err := row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, t)
}

func postTodoHandler(c *gin.Context) {
	db := dbConnect()
	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	query := `
	INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id;
	`
	var id int
	row := db.QueryRow(query, t.Title, t.Status)
	err := row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "Success",
		"id":     id,
	})
}

func deleteTodoByIdHandler(c *gin.Context) {
	db := dbConnect()
	id, _ := strconv.Atoi(c.Param("id"))
	stmt, err := db.Prepare("DELETE FROM todos WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if _, err := stmt.Exec(id); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
	})
}
