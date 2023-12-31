package main

import (
	"database/sql"
	"log"

	sqlMigrations "github.com/jonfriesen/todo-templ-htmx-tailwind/sql"
	goose "github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"

	sqlcdb "github.com/jonfriesen/todo-templ-htmx-tailwind/internal/db"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/service"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/todo"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrations := sqlMigrations.EmbeddedFiles
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatalf("setting database dialect: %v", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("migrating database: %v", err)
	}

	querier := sqlcdb.New(db)

	app := service.Application{
		TodoService: todo.New(querier),
	}

	app.ServeHTTP()
}
