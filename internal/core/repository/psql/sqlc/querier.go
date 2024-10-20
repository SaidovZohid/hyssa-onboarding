// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"context"
)

type Querier interface {
	CountTodo(ctx context.Context) (int64, error)
	CreateTodo(ctx context.Context, title string) (Todo, error)
	GetTodoById(ctx context.Context, id int32) (Todo, error)
	ListTodo(ctx context.Context, arg ListTodoParams) ([]Todo, error)
	PatchTodo(ctx context.Context, arg PatchTodoParams) (Todo, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (Todo, error)
}

var _ Querier = (*Queries)(nil)
