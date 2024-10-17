package service

import (
	"context"
	"encoding/json"

	"learning/hyssa-learn/generated/todo_service"
	"learning/hyssa-learn/internal/core/repository"
	"learning/hyssa-learn/internal/core/repository/psql/sqlc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoService struct {
	todo_service.UnimplementedTodoServiceServer
	store repository.Store
}

func NewTodoService(store repository.Store) *TodoService {
	return &TodoService{
		store: store,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, req *todo_service.Todo) (*todo_service.Todo, error) {
	data, err := s.store.CreateTodo(ctx, req.Title)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err.Error())
	}

	response := todo_service.Todo{}
	if todoInBytes, err := json.Marshal(data); err == nil {
		json.Unmarshal(todoInBytes, &response)
	} else {
		return nil, status.Errorf(codes.Internal, "failed to marshal todo: %s    ", err.Error())
	}
	return &response, nil
}

func (s *TodoService) UpdateTodoById(ctx context.Context, req *todo_service.Todo) (*todo_service.Todo, error) {
	data, err := s.store.PatchTodo(ctx, sqlc.PatchTodoParams{
		ID:        req.Id,
		Title:     req.Title,
		Completed: req.Completed,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to patch user: %s", err.Error())
	}

	response := todo_service.Todo{}
	if todoInBytes, err := json.Marshal(data); err == nil {
		json.Unmarshal(todoInBytes, &response)
	} else {
		return nil, status.Errorf(codes.Internal, "failed to marshal todo: %s    ", err.Error())
	}
	return &response, nil
}

func (s *TodoService) GetAllTodos(ctx context.Context, req *todo_service.GetAllTodosReq) (*todo_service.GetAllTodosResponse, error) {
	offset := (req.Page - 1) * req.Limit
	todos, err := s.store.ListTodo(ctx, sqlc.ListTodoParams{
		Limit:  req.Limit,
		Offset: offset,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todos: %s", err.Error())
	}
	response := todo_service.GetAllTodosResponse{}
	for _, todo := range todos {
		todoInBytes, err := json.Marshal(todo)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to marshal todo: %s    ", err.Error())
		}
		var todoResponse todo_service.Todo
		json.Unmarshal(todoInBytes, &todoResponse)
		response.Todos = append(response.Todos, &todoResponse)
	}
	return &response, nil
}

func (s *TodoService) GetTodoById(ctx context.Context, req *todo_service.TodoPk) (*todo_service.Todo, error) {
	todo, err := s.store.GetTodoById(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todo: %s", err.Error())
	}

	response := todo_service.Todo{}
	if todoInBytes, err := json.Marshal(todo); err == nil {
		json.Unmarshal(todoInBytes, &response)
	} else {
		return nil, status.Errorf(codes.Internal, "failed to marshal todo: %s    ", err.Error())
	}

	return &response, nil
}
