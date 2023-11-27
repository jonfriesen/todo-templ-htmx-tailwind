package service

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/components"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/todo"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	comp := components.Page(app.TodoService.GetItems())

	comp.Render(r.Context(), w)
}

func (app *Application) completeItem(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "item_id")
	fmt.Println("toggling todo", todoID)
	fmt.Printf("%+v\n", app.TodoService.GetItems())

	todo, err := app.TodoService.ToggleComplete(todoID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Println(todo)

	todoRow := components.TodoRow(todo)

	todoRow.Render(r.Context(), w)
}

type newTodo struct {
	Description string `form:"description"`
}

func (app *Application) addItem(w http.ResponseWriter, r *http.Request) {
	var newTodo newTodo
	err := DecodePostForm(r, &newTodo)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	fmt.Printf("form submission: %+v\n", newTodo)

	if newTodo.Description == "" {
		app.badRequest(w, r, fmt.Errorf("todo item description required"))
		return
	}

	todo := app.TodoService.AddItem(todo.TodoItem{Description: newTodo.Description})

	todoRow := components.TodoRow(todo)

	todoRow.Render(r.Context(), w)
}
