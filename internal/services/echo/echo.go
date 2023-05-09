package echo

import (
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Healthcheck(context.Context) error
	Echo(ctx context.Context, message string) (string, error)
}

var _ Service = (*service)(nil)

type service struct {
	log *zap.Logger
}

func (s *service) Healthcheck(context.Context) error {
	// dependencies healthcheck
	s.log.Info("echo service healthcheck")
	return nil
}

func (s *service) Echo(_ context.Context, message string) (string, error) {
	s.log.Info("echo service echo", zap.String("message", message))
	return message, nil
}

func New(log *zap.Logger) Service {
	return &service{
		log: log,
	}
}
