package distributor

import (
	"context"
	"encoding/json"
	"errors"

	"learning/hyssa-learn/internal/transport/queue/processor"

	"github.com/hibiken/asynq"

	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opt ...asynq.Option,
) error {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.New("failed to marshal payload: " + err.Error())
	}

	for taskName := range processor.Processors {
		task := asynq.NewTask(taskName, jsonBytes, opt...)
		info, err := distributor.client.EnqueueContext(ctx, task)
		if err != nil {
			return errors.New("failed to enqueue task: " + err.Error())
		}
		log.Info().Str("type:", task.Type()).Bytes("paylaod", jsonBytes).
			Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	}

	return nil
}
