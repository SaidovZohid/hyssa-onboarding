package grpc

import (
	"learning/hyssa-learn/generated/todo_service"
	"learning/hyssa-learn/internal/config"
	"learning/hyssa-learn/internal/core/repository"
	"learning/hyssa-learn/internal/core/service"
	"learning/hyssa-learn/internal/transport/grpc/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func New(repo repository.Store, cfg *config.Config) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middleware.GrpcLoggerMiddleware,
			),
		),
	)

	reflection.Register(grpcServer)
	todo_service.RegisterTodoServiceServer(grpcServer, service.NewTodoService(repo))
	return grpcServer
}
