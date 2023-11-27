package todo

import (
	"fmt"
	"sync"

	"github.com/rs/xid"
)

// TodoItem represents a single todo item.
type TodoItem struct {
	ID          string
	Description string
	Complete    bool
}

// TodoList holds a list of todo items.
type TodoList struct {
	items map[string]TodoItem
	mu    sync.Mutex
}

func New() TodoList {
	return TodoList{
		items: make(map[string]TodoItem),
	}
}

// AddItem adds a new item to the TodoList.
func (t *TodoList) AddItem(item TodoItem) TodoItem {
	t.mu.Lock()
	defer t.mu.Unlock()
	item.ID = xid.New().String()
	t.items[item.ID] = item

	return item
}

// GetItems returns a copy of the list of items.
func (t *TodoList) GetItems() []TodoItem {
	t.mu.Lock()
	defer t.mu.Unlock()
	resp := []TodoItem{}
	for _, i := range t.items {
		resp = append(resp, i)
	}
	return resp
}

// ToggleComplete toggles the completion status of an item.
func (t *TodoList) ToggleComplete(id string) (TodoItem, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if item, ok := t.items[id]; ok {
		item.Complete = !item.Complete
		t.items[item.ID] = item

		return item, nil
	}

	return TodoItem{}, fmt.Errorf("not found")
}
