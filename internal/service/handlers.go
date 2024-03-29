package service

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/components"
	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/db"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	components.Index().Render(r.Context(), w)
}

func (app *Application) todo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	items, err := app.TodoService.GetItems(r.Context(), userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user, err := app.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	comp := components.Page(components.TodoPage(user.Name, items))

	comp.Render(r.Context(), w)
}

func (app *Application) completeItem(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "item_id")
	fmt.Println("toggling todo", todoID)

	todo, err := app.TodoService.ToggleComplete(r.Context(), todoID)
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
	userID := r.Context().Value("user_id").(string)

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

	todo, err := app.TodoService.AddItem(r.Context(), &db.TodoItem{UserID: userID, Description: newTodo.Description})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	todoRow := components.TodoRow(todo)

	todoRow.Render(r.Context(), w)
}
