package server

import "context"

type Server struct{}

func (s *Server) GetTodos(ctx context.Context, req GetTodosRequestObject) (GetTodosResponseObject, error) {
	staticTodo := Todo{
		Id:    1,
		Title: "Implement dynamic todos",
	}

	return GetTodos200JSONResponse{Items: []Todo{staticTodo}}, nil
}

func (s *Server) PostTodos(ctx context.Context, req PostTodosRequestObject) (PostTodosResponseObject, error) {
	return PostTodos201JSONResponse{}, nil
}
