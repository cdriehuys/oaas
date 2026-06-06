package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/cdriehuys/oaas/api/internal/repositories"
)

type todoRepo interface {
	CreateTodo(ctx context.Context, title string) (repositories.Todo, error)
	GetComplete(ctx context.Context) ([]repositories.Todo, error)
	GetIncomplete(ctx context.Context) ([]repositories.Todo, error)
	SetComplete(ctx context.Context, id int) (repositories.Todo, error)
	SetIncomplete(ctx context.Context, id int) (repositories.Todo, error)
}

type Server struct {
	todos todoRepo
}

func NewServer(todos todoRepo) *Server {
	return &Server{todos}
}

func (s *Server) GetTodos(ctx context.Context, req GetTodosRequestObject) (GetTodosResponseObject, error) {
	todos, err := s.todos.GetIncomplete(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving todos: %v", err)
	}

	return GetTodos200JSONResponse(todoCollectionToRep(todos)), nil
}

func (s *Server) PostTodos(ctx context.Context, req PostTodosRequestObject) (PostTodosResponseObject, error) {
	created, err := s.todos.CreateTodo(ctx, req.Body.Title)
	if err != nil {
		return nil, fmt.Errorf("creating todo: %v", err)
	}

	return PostTodos201JSONResponse(todoToRep(created)), nil
}

func (s *Server) PutTodosTodoIDState(ctx context.Context, req PutTodosTodoIDStateRequestObject) (PutTodosTodoIDStateResponseObject, error) {
	if !req.Body.Valid() {
		return PutTodosTodoIDState400JSONResponse{BadRequestJSONResponse{Message: fmt.Sprintf("Invalid state %q.", *req.Body)}}, nil
	}

	var todo repositories.Todo
	var err error
	if *req.Body == Complete {
		todo, err = s.todos.SetComplete(ctx, req.TodoID)
	} else {
		todo, err = s.todos.SetIncomplete(ctx, req.TodoID)
	}

	if err != nil {
		if errors.Is(err, repositories.ErrTodoNotFound) {
			return PutTodosTodoIDState404JSONResponse{TodoNotFoundJSONResponse{Message: fmt.Sprintf("No todo with ID %d.", req.TodoID)}}, nil
		}

		return nil, fmt.Errorf("setting todo state: %v", err)
	}

	return PutTodosTodoIDState200JSONResponse(todoToRep(todo)), nil
}

func todoToRep(todo repositories.Todo) Todo {
	return Todo{
		Id:          todo.ID,
		Title:       todo.Title,
		CompletedAt: todo.CompletedAt,
	}
}

func todoCollectionToRep(todos []repositories.Todo) TodoCollection {
	reps := make([]Todo, len(todos))
	for i, todo := range todos {
		reps[i] = todoToRep(todo)
	}

	return TodoCollection{Items: reps}
}
