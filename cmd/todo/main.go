package main

import (
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/service"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/todo"
)

func main() {
	app := service.Application{
		TodoService: todo.New(),
	}

	// seed with some sample todos
	app.TodoService.AddItem(todo.TodoItem{
		Description: "Walk the dogs",
		Complete:    false,
	})
	app.TodoService.AddItem(todo.TodoItem{
		Description: "File Taxes",
		Complete:    true,
	})

	app.ServeHTTP()

}
