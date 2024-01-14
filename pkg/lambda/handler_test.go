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

func TestHandler_EventHandler_Sync(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		tests  = []struct {
			name    string
			give    json.RawMessage
			want    any
			wantErr assert.ErrorAssertionFunc
		}{
			{
				name:    "should handle an API Gateway event",
				give:    testutils.ReadJSONFromFile(t, "../../testutils/data/apigw.json"),
				want:    &service.Response{Status: "ok"},
				wantErr: assert.NoError,
			},
			{
				name:    "should return error if cannot handle an API Gateway event",
				give:    testutils.ReadJSONFromFile(t, "../../testutils/data/bad_apigw.json"),
				wantErr: assert.Error,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler[service.Request, service.Response](service.New().Do)
			r, err := h.EventHandler(ctx, tt.give)

			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, r)
			}
		})
	}
}

func TestHandler_EventHandler_Async(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		tests  = []struct {
			name    string
			give    json.RawMessage
			wantErr assert.ErrorAssertionFunc
		}{
			{
				name:    "should handle a SQS event",
				give:    testutils.ReadJSONFromFile(t, "../../testutils/data/sqs.json"),
				wantErr: assert.Error,
			},
			{
				name:    "should handle a SNS event",
				give:    testutils.ReadJSONFromFile(t, "../../testutils/data/sns.json"),
				wantErr: assert.NoError,
			},
			{
				name:    "should handle a DynamoDB event",
				give:    testutils.ReadJSONFromFile(t, "../../testutils/data/dynamodb.json"),
				wantErr: assert.Error,
			},
			{
				name:    "should handle a CloudWatch event",
				give:    testutils.ReadJSONFromFile(t, "../../testutils/data/cloudwatch.json"),
				wantErr: assert.NoError,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(service.New().Do)
			r, err := h.EventHandler(ctx, tt.give)

			tt.wantErr(t, err)
			assert.Equal(t, nil, r)
		})
	}
}

func TestHandler_EventHandler_Batch(t *testing.T) {
	var (
		ctx, _ = logger.NewContextWithLogger(context.TODO())
		tests  = []struct {
			name    string
			give    json.RawMessage
			want    any
			wantErr assert.ErrorAssertionFunc
		}{
			{
				name: "should handle a batch SQS event",
				give: testutils.ReadJSONFromFile(t, "../../testutils/data/sqs.json"),
				want: events.SQSEventResponse{BatchItemFailures: []events.SQSBatchItemFailure{
					{
						ItemIdentifier: "MessageID_2",
					},
				}},
				wantErr: assert.NoError,
			},
			{
				name: "should handle a batch DynamoDB event",
				give: testutils.ReadJSONFromFile(t, "../../testutils/data/dynamodb.json"),
				want: events.DynamoDBEventResponse{BatchItemFailures: []events.DynamoDBBatchItemFailure{
					{
						ItemIdentifier: "f07f8ca4b0b26cb9c4e5e77e42f274ee",
					},
				}},
				wantErr: assert.NoError,
			},
			{
				name:    "should not handle an invalid event",
				give:    testutils.ReadJSONFromFile(t, "../../testutils/data/cloudwatch.json"),
				wantErr: assert.Error,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(service.New().Do)
			h.Batch()
			r, err := h.EventHandler(ctx, tt.give)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, r)
		})
	}
}
