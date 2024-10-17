package handlers

import (
	"context"
	"fmt"

	"learning/hyssa-learn/generated/todo_service"
	"learning/hyssa-learn/internal/config"
	"learning/hyssa-learn/pkg/wrapper"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcMain "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(ctx context.Context, cfg *config.Config) *runtime.ServeMux {
	gwMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(wrapper.CustomMatcher),
	)
	connPingService, err := grpcMain.Dial(
		fmt.Sprintf("%s:%d", cfg.Grpc.Host, cfg.Grpc.Port),
		grpcMain.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil
	}

	if err := todo_service.RegisterTodoServiceHandler(ctx, gwMux, connPingService); err != nil {
		return nil
	}

	return gwMux
}

func makeHost(host string, port int32) string {
	return host + ":" + fmt.Sprintf("%d", port)
}
