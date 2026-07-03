package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/cdriehuys/oaas/api/internal/repositories"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const probTypeValidationError string = "/probs/validation-error"

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
	showComplete := *req.Params.State == GetTodosParamsStateComplete

	var todos []repositories.Todo
	var err error
	if showComplete {
		todos, err = s.todos.GetComplete(ctx)
	} else {
		todos, err = s.todos.GetIncomplete(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("retrieving todos: %v", err)
	}

	return GetTodos200JSONResponse(todoCollectionToRep(todos)), nil
}

func (t *PostTodosJSONRequestBody) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Title, validation.Required, validation.Length(1, 100)),
	)
}

func (s *Server) PostTodos(ctx context.Context, req PostTodosRequestObject) (PostTodosResponseObject, error) {
	if err := req.Body.Validate(); err != nil {
		if resp, ok := asValidationError(err); ok {
			return PostTodos400ApplicationProblemPlusJSONResponse{resp}, nil
		}

		return nil, fmt.Errorf("validating todo: %v", err)
	}

	created, err := s.todos.CreateTodo(ctx, req.Body.Title)
	if err != nil {
		return nil, fmt.Errorf("creating todo: %v", err)
	}

	return PostTodos201JSONResponse(todoToRep(created)), nil
}

func (s *Server) PutTodosTodoIDState(ctx context.Context, req PutTodosTodoIDStateRequestObject) (PutTodosTodoIDStateResponseObject, error) {
	if !req.Body.Valid() {
		errType := "/probs/invalid-state"
		title := "Invalid State"
		detail := fmt.Sprintf("The state %q is invalid.", *req.Body)

		return PutTodosTodoIDState400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse{
				Type:   &errType,
				Title:  &title,
				Detail: &detail,
			},
		}, nil
	}

	var todo repositories.Todo
	var err error
	if *req.Body == TodoStateComplete {
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

func asValidationError(err error) (BadRequestApplicationProblemPlusJSONResponse, bool) {
	if vErr, ok := errors.AsType[validation.Errors](err); ok {
		return validationErrorToResponse(vErr), true
	}

	return BadRequestApplicationProblemPlusJSONResponse{}, false
}

func validationErrorToResponse(errs validation.Errors) BadRequestApplicationProblemPlusJSONResponse {
	var fieldErrors []FieldError
	for field, err := range errs {
		code := err.(validation.Error).Code()
		title := err.Error()

		fieldErrors = append(fieldErrors, FieldError{code, field, title})
	}

	errType := probTypeValidationError
	title := "Bad Request"

	return BadRequestApplicationProblemPlusJSONResponse{
		Type:        &errType,
		Title:       &title,
		FieldErrors: &fieldErrors,
	}
}
