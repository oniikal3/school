package todo

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/oniikal3/school/database"

	"github.com/gin-gonic/gin"
)

func main() {}

type Todo struct {
	ID     int    `json: "id"`
	Title  string `json: "title"`
	Status string `json: "status"`
}

func dbConnect() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	return db
}

func GetHandler(c *gin.Context) {
	db := dbConnect()
	id, _ := strconv.Atoi(c.Param("id"))
	row, err := database.GetTodo(db, id)
	t := Todo{}
	row.Scan(&t.ID, &t.Title, &t.Status)

	defer db.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, t)
}

func PostHandler(c *gin.Context) {
	db := dbConnect()
	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id, err := database.PostTodo(db, t.Title, t.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "Success",
		"id":     id,
	})
}

func DeleteHandler(c *gin.Context) {
	db := dbConnect()
	id, _ := strconv.Atoi(c.Param("id"))
	err := database.DeleteTodo(db, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "Success",
	})
}
