package repositories

import (
	"errors"
	"time"
)

var ErrTodoNotFound = errors.New("todo not found")

type Todo struct {
	ID          int
	Title       string
	CompletedAt *time.Time
}
