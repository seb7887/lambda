package lambda

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"seb7887/lambda/internal/service"
	"seb7887/lambda/pkg/logger"
	"seb7887/lambda/testutils"
	"testing"
)

func TestHandler_SNSEventHandler(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		h      = New(service.New())
	)

	assert.NoError(t, h.SNSEventHandler(ctx, testutils.SNSMockEvent))
}

func TestHandler_SQSEventHandler(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		h      = New(service.New())
	)

	assert.NoError(t, h.SQSEventHandler(ctx, testutils.SQSMockEvent))
}

func TestHandler_BatchSQSEventHandler(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		h      = New(service.New())
	)

	r, err := h.BatchSQSEventHandler(ctx, testutils.BatchSQSMockEvent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.BatchItemFailures))
	assert.Equal(t, "2", r.BatchItemFailures[0].ItemIdentifier)
}

func TestHandler_DDBEventHandler(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		h      = New(service.New())
		evt    = prepareEvent[events.DynamoDBEvent](t, "../../testutils/ddbevent.json")
	)

	assert.Error(t, h.DDBEventHandler(ctx, evt))
}

func TestHandler_BatchDDBEventHandler(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		h      = New(service.New())
		evt    = prepareEvent[events.DynamoDBEvent](t, "../../testutils/ddbevent.json")
	)

	r, err := h.BatchDDBEventHandler(ctx, evt)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r.BatchItemFailures))
	assert.Equal(t, "f07f8ca4b0b26cb9c4e5e77e42f274ee", r.BatchItemFailures[0].ItemIdentifier)
}

func prepareEvent[T any](t *testing.T, filename string) T {
	inputJSON := testutils.ReadJSONFromFile(t, filename)
	var evt T
	if err := json.Unmarshal(inputJSON, &evt); err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}
	return evt
}
