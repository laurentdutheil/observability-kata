package repository

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"todo_odd/domain"
)

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository() *SqliteRepository {
	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := "create table todo (id integer not null primary key, title text, description text);"

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	return &SqliteRepository{db}
}

func (r SqliteRepository) AddTodo(ctx context.Context, title string, description string) domain.Todo {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into todo(id, title, description) values(?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	result, err := stmt.Exec(nil, title, description)
	if err != nil {
		log.Fatal(err)
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return domain.Todo{
		Id:          int(lastInsertId),
		Title:       title,
		Description: description,
	}
}

func (r SqliteRepository) Get(id int) (domain.Todo, error) {
	stmt, err := r.db.Prepare("select id, title, description from todo where id=?")
	if err != nil {
		log.Fatal(err)
		return domain.Todo{}, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	var qId int
	var qTitle string
	var qDescription string
	err = stmt.QueryRow(id).Scan(&qId, &qTitle, &qDescription)
	if err != nil {
		log.Println(err)
		return domain.Todo{}, err
	}

	return domain.Todo{
		Id:          qId,
		Title:       qTitle,
		Description: qDescription,
	}, nil
}

func (r SqliteRepository) All() []domain.Todo {
	rows, err := r.db.Query("select id, title, description from todo")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var allTodos []domain.Todo
	for rows.Next() {
		var qId int
		var qTitle string
		var qDescription string
		err = rows.Scan(&qId, &qTitle, &qDescription)
		if err != nil {
			log.Fatal(err)
		}
		allTodos = append(allTodos, domain.Todo{
			Id:          qId,
			Title:       qTitle,
			Description: qDescription,
		})
	}

	return allTodos
}

func (r SqliteRepository) Close() {
	_ = r.db.Close()
}
