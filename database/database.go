package database

import (
	"database/sql"
)

func main() {}

func GetTodo(db *sql.DB, id int) (*sql.Row, error) {
	stmt, err := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1;")
	row := stmt.QueryRow(id)
	return row, err
}

func PostTodo(db *sql.DB, title string, status string) (int, error) {
	query := `
	INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id;
	`
	var id int
	row := db.QueryRow(query, title, status)
	err := row.Scan(&id)
	return id, err
}

func DeleteTodo(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM todos WHERE id=$1;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}
