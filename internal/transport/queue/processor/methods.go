package processor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
)

const (
	TaskSendVerifyEmail     = "task:send_verify_email-alif"
	TaskSendVerifyEmailIman = "email:send_verify_email-iman"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (process *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return errors.New("failed to unmarshal payload: " + asynq.SkipRetry.Error())
	}

	// user, err := process.store.GetUserByName(ctx, payload.Username)
	// if err != nil {
	// 	if err == pgx.ErrNoRows {
	// 		return errors.New("user doesn't exists: " + err.Error())
	// 	}
	// 	return errors.New("failed to get user: " + err.Error())
	// }
	// TODO: send email to user logic
	// logger.Log.Info(fmt.Sprintf("type: %s , payload: %s , email: %s , processed task", task.Type(), string(task.Payload()), user.Email))
	return nil
}

func (process *RedisTaskProcessor) ProcessTaskSendVerifyEmail2(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return errors.New("failed to unmarshal payload: " + asynq.SkipRetry.Error())
	}

	// user, err := process.store.GetUserByName(ctx, payload.Username)
	// if err != nil {
	// 	if err == pgx.ErrNoRows {
	// 		return errors.New("user doesn't exists: " + err.Error())
	// 	}
	// 	return errors.New("failed to get user: " + err.Error())
	// }
	// TODO: send email to user logic
	// logger.Log.Info(fmt.Sprintf("type: %s , payload: %s , email: %s ,processed task for second one", task.Type(), string(task.Payload()), user.Email))
	return nil
}
