package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type HandlerFn[I, O any] func(ctx context.Context, in *I) (*O, error)

type Handler[I, O any] struct {
	isBatch bool
	fn      HandlerFn[I, O]
}

func NewHandler[I, O any](fn HandlerFn[I, O], middlewares ...MiddlewareFn[I, O]) *Handler[I, O] {
	// Apply middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		fn = middlewares[i](fn)
	}

	return &Handler[I, O]{
		fn: fn,
	}
}

func (h *Handler[I, O]) Batch() {
	h.isBatch = true
}

func (h *Handler[I, O]) EventHandler(ctx context.Context, raw json.RawMessage) (any, error) {
	if h.isBatch {
		return h.batch(ctx, raw)
	}

	switch eventSource(raw) {
	case _apiEvent:
		res, err := h.withResponse(ctx, raw)
		return res, err
	default:
		return nil, h.withoutResponse(ctx, raw)
	}
}

func (h *Handler[I, O]) withResponse(ctx context.Context, raw json.RawMessage) (*O, error) {
	in, err := adaptJSON[I](raw)
	if err != nil {
		return nil, err
	}

	return h.fn(ctx, in)
}

func (h *Handler[I, O]) withoutResponse(ctx context.Context, raw json.RawMessage) error {
	evt, err := identifyEventType(raw)
	if err != nil {
		return err
	}

	return h.handle(ctx, evt)
}

func (h *Handler[I, O]) handle(ctx context.Context, evt any) error {
	var (
		processEvent = func(data []byte) error {
			var in I
			if err := json.Unmarshal(data, &in); err != nil {
				return err
			}

			_, err := h.fn(ctx, &in)
			return err
		}
	)

	switch event := evt.(type) {
	case *events.SQSEvent:
		for _, record := range event.Records {
			if err := processEvent([]byte(record.Body)); err != nil {
				return err
			}
		}
	case *events.SNSEvent:
		for _, record := range event.Records {
			if err := processEvent([]byte(record.SNS.Message)); err != nil {
				return err
			}
		}
	case *events.DynamoDBEvent:
		for _, record := range event.Records {
			in, err := h.unmarshalStreamImage(record.Change.NewImage)
			if err != nil {
				return err
			}

			_, err = h.fn(ctx, &in)
			if err != nil {
				return err
			}
		}
	case *events.CloudWatchEvent:
		if err := processEvent(event.Detail); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler[I, O]) batch(ctx context.Context, raw json.RawMessage) (any, error) {
	evt, err := identifyEventType(raw)
	if err != nil {
		return nil, err
	}

	switch event := evt.(type) {
	case *events.SQSEvent:
		return h.processSQSBatch(ctx, event)
	case *events.DynamoDBEvent:
		return h.processDynamoDBBatch(ctx, event)
	default:
		return nil, errors.New("invalid batch event")
	}
}

func (h *Handler[I, O]) unmarshalStreamImage(attr map[string]events.DynamoDBAttributeValue) (I, error) {
	var (
		out     I
		attrMap = make(map[string]*dynamodb.AttributeValue)
	)

	for k, v := range attr {
		b, err := v.MarshalJSON()
		if err != nil {
			return out, err
		}
		var dbAttr dynamodb.AttributeValue
		if err = json.Unmarshal(b, &dbAttr); err != nil {
			return out, err
		}
		attrMap[k] = &dbAttr
	}

	if err := dynamodbattribute.UnmarshalMap(attrMap, &out); err != nil {
		return out, err
	}

	return out, nil
}

func (h *Handler[I, O]) processEvent(ctx context.Context, body string) (*O, error) {
	var in I
	if err := json.Unmarshal([]byte(body), &in); err != nil {
		return nil, err
	}
	return h.fn(ctx, &in)
}

func (h *Handler[I, O]) processSQSBatch(ctx context.Context, event *events.SQSEvent) (events.SQSEventResponse, error) {
	var failedMessages []events.SQSBatchItemFailure
	for _, record := range event.Records {
		if _, err := h.processEvent(ctx, record.Body); err != nil {
			failedMessages = append(failedMessages, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
		}
	}
	return events.SQSEventResponse{BatchItemFailures: failedMessages}, nil
}

func (h *Handler[I, O]) processDynamoDBBatch(ctx context.Context, event *events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	var failedItems []events.DynamoDBBatchItemFailure
	for _, record := range event.Records {
		in, err := h.unmarshalStreamImage(record.Change.NewImage)
		if err != nil {
			failedItems = append(failedItems, events.DynamoDBBatchItemFailure{ItemIdentifier: record.EventID})
			continue
		}

		if _, err = h.fn(ctx, &in); err != nil {
			failedItems = append(failedItems, events.DynamoDBBatchItemFailure{ItemIdentifier: record.EventID})
		}
	}
	return events.DynamoDBEventResponse{BatchItemFailures: failedItems}, nil
}
