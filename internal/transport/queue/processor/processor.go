package processor

import (
	"context"

	"learning/hyssa-learn/internal/core/repository"

	"github.com/hibiken/asynq"
)

type ProcessorFunction func(ctx context.Context, task *asynq.Task) error

var Processors = map[string]ProcessorFunction{}

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  repository.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store repository.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{},
	)
	processor := RedisTaskProcessor{
		store:  store,
		server: server,
	}

	Processors = map[string]ProcessorFunction{
		TaskSendVerifyEmail:     processor.ProcessTaskSendVerifyEmail,
		TaskSendVerifyEmailIman: processor.ProcessTaskSendVerifyEmail2,
	}

	return &processor
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	for taskName, processor := range Processors {
		mux.HandleFunc(taskName, processor)
	}

	return processor.server.Start(mux)
}
