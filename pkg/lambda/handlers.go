package lambda

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"
	"seb7887/lambda/internal/service"
	"seb7887/lambda/pkg/logger"
)

type HandlerService interface {
	Do(ctx context.Context, in *service.Request) (*service.Response, error)
}

type Handlers struct {
	svc HandlerService
}

func New(svc HandlerService) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) APIEventHandler(ctx context.Context, in *service.Request) (*service.Response, error) {
	ctx, log := logger.NewContextWithLogger(ctx)
	log.With(zap.Any("event", in)).Info("start")
	res, err := h.svc.Do(ctx, in)
	if err != nil {
		log.With(zap.Any("error", err)).Error("execution error")
		return nil, err
	}

	log.With(zap.Any("response", res)).Info("success")
	return res, nil
}

func (h *Handlers) SNSEventHandler(ctx context.Context, evt events.SNSEvent) error {
	ctx, log := logger.NewContextWithLogger(ctx)

	for _, record := range evt.Records {
		var (
			in service.Request
		)

		if err := json.Unmarshal([]byte(record.SNS.Message), &in); err != nil {
			log.With(zap.Any("event.raw", record)).Error("cannot deserialize event")
			return err
		}

		log.With(zap.Any("event", in)).Info("start")
		_, err := h.svc.Do(ctx, &in)
		if err != nil {
			log.With(zap.Any("error", err)).Error("execution error")
			return err
		}

		log.Info("success")
	}

	return nil
}

func (h *Handlers) SQSEventHandler(ctx context.Context, evt events.SQSEvent) error {
	ctx, log := logger.NewContextWithLogger(ctx)

	for _, record := range evt.Records {
		var in service.Request

		if err := json.Unmarshal([]byte(record.Body), &in); err != nil {
			log.With(zap.Any("event.raw", record)).Error("cannot deserialize event")
			return err
		}

		log.With(zap.Any("event", in)).Info("start")
		_, err := h.svc.Do(ctx, &in)
		if err != nil {
			log.With(zap.Any("error", err)).Error("execution error")
			return err
		}

		log.Info("success")
	}

	return nil
}

func (h *Handlers) BatchSQSEventHandler(ctx context.Context, evt events.SQSEvent) (events.SQSEventResponse, error) {
	var failedMessages []events.SQSBatchItemFailure
	ctx, log := logger.NewContextWithLogger(ctx)

	for _, record := range evt.Records {
		var in service.Request

		if err := json.Unmarshal([]byte(record.Body), &in); err != nil {
			log.With(zap.Any("event.raw", record)).Error("cannot deserialize event")
			failedMessages = append(failedMessages, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			continue
		}

		log.With(zap.Any("event", in)).Info("start")
		_, err := h.svc.Do(ctx, &in)
		if err != nil {
			log.With(zap.Any("error", err)).Error("execution error")
			failedMessages = append(failedMessages, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			continue
		}

		log.Info("success")
	}

	return events.SQSEventResponse{BatchItemFailures: failedMessages}, nil
}

func (h *Handlers) DDBEventHandler(ctx context.Context, evt events.DynamoDBEvent) error {
	ctx, log := logger.NewContextWithLogger(ctx)

	for _, record := range evt.Records {
		var in service.Request
		if err := h.unmarshalStreamImage(record.Change.NewImage, &in); err != nil {
			log.With(zap.Any("event.raw", record), zap.Any("error", err)).Error("cannot deserialize event")
			return err
		}

		log.With(zap.Any("event", in)).Info("start")
		_, err := h.svc.Do(ctx, &in)
		if err != nil {
			log.With(zap.Any("error", err)).Error("execution error")
			return err
		}

		log.Info("success")
	}

	return nil
}

func (h *Handlers) BatchDDBEventHandler(ctx context.Context, evt events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	var failedItems []events.DynamoDBBatchItemFailure
	ctx, log := logger.NewContextWithLogger(ctx)

	for _, record := range evt.Records {
		var in service.Request
		if err := h.unmarshalStreamImage(record.Change.NewImage, &in); err != nil {
			log.With(zap.Any("event.raw", record), zap.Any("error", err)).Error("cannot deserialize event")
			failedItems = append(failedItems, events.DynamoDBBatchItemFailure{ItemIdentifier: record.EventID})
			continue
		}

		log.With(zap.Any("event", in)).Info("start")
		_, err := h.svc.Do(ctx, &in)
		if err != nil {
			log.With(zap.Any("error", err)).Error("execution error")
			failedItems = append(failedItems, events.DynamoDBBatchItemFailure{ItemIdentifier: record.EventID})
			continue
		}

		log.Info("success")
	}

	return events.DynamoDBEventResponse{
		BatchItemFailures: failedItems,
	}, nil
}

func (h *Handlers) unmarshalStreamImage(attr map[string]events.DynamoDBAttributeValue, out any) error {
	attrMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attr {
		b, err := v.MarshalJSON()
		if err != nil {
			return err
		}
		var dbAttr dynamodb.AttributeValue
		if err = json.Unmarshal(b, &dbAttr); err != nil {
			return err
		}
		attrMap[k] = &dbAttr
	}

	return dynamodbattribute.UnmarshalMap(attrMap, out)
}

func (h *Handlers) MultiEventHandler(ctx context.Context, raw json.RawMessage) (any, error) {
	source := eventSource(raw)
	ctx, log := logger.NewContextWithLogger(ctx)
	log.With(zap.Any("event.source", source)).Info("event parsed")

	if source == _apiEvent {
		return h.api(ctx, raw), nil
	}

	if source == _sqsEvent {
		return nil, h.sqs(ctx, raw)
	}

	return nil, nil
}

func (h *Handlers) api(ctx context.Context, raw json.RawMessage) *service.Response {
	var in service.Request
	if err := json.Unmarshal(raw, &in); err != nil {
		return nil
	}
	r, err := h.svc.Do(ctx, &in)
	if err != nil {
		return nil
	}
	return r
}

func (h *Handlers) sqs(ctx context.Context, raw json.RawMessage) error {
	var evt events.SQSEvent
	if err := json.Unmarshal(raw, &evt); err != nil {
		return err
	}

	for _, record := range evt.Records {
		var in service.Request
		if err := json.Unmarshal([]byte(record.Body), &in); err != nil {
			return err
		}

		_, err := h.svc.Do(ctx, &in)
		if err != nil {
			return err
		}
	}

	return nil
}
