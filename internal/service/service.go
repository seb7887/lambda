package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"seb7887/lambda/pkg/logger"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

type Request struct {
	UserID      string `json:"user_uuid"`
	OperationID string `json:"operation_id"`
}

type Response struct {
	Status string `json:"status"`
}

func (s *Service) Do(ctx context.Context, in *Request) (*Response, error) {
	log := logger.New(ctx)
	log.With(zap.Any("request", in)).Debug("service")
	if in.UserID == "x" {
		return nil, errors.New("oops")
	}

	return &Response{Status: "ok"}, nil
}
