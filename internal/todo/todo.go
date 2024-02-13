package todo

import (
	"context"

	"github.com/jonfriesen/todo-templ-htmx-tailwind/internal/db"
	"github.com/rs/xid"
)

// TodoList holds a list of todo items.
type TodoList struct {
	db db.Querier
}

func New(db db.Querier) TodoList {
	return TodoList{
		db: db,
	}
}

// AddItem adds a new item to the TodoList.
func (t *TodoList) AddItem(ctx context.Context, item *db.TodoItem) (*db.TodoItem, error) {
	item.ID = xid.New().String()
	t.db.InsertTodo(ctx, &db.InsertTodoParams{
		ID:          item.ID,
		UserID:      item.UserID,
		Description: item.Description,
		Complete:    item.Complete,
	})

	var err error
	item, err = t.db.GetTodo(ctx, item.ID)

	return item, err
}

// GetItems returns a copy of the list of items.
func (t *TodoList) GetItems(ctx context.Context, userID string) ([]*db.TodoItem, error) {
	return t.db.ListTodos(ctx, userID)
}

// ToggleComplete toggles the completion status of an item.
func (t *TodoList) ToggleComplete(ctx context.Context, id string) (*db.TodoItem, error) {
	err := t.db.CompleteTodo(ctx, id)
	if err != nil {
		return nil, err
	}

	return t.db.GetTodo(ctx, id)
}
