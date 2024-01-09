package testutils

import (
	"github.com/aws/aws-lambda-go/events"
)

var (
	Payload    = `{"user_uuid":"iao","operation_id":"123"}`
	BadPayload = `{"user_uuid":"x","operation_id":"123"}`

	SNSMockEvent = events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				EventSource: "aws:sns",
				SNS: events.SNSEntity{
					Message: Payload,
				},
			},
			{
				EventSource: "aws:sns",
				SNS: events.SNSEntity{
					Message: Payload,
				},
			},
		},
	}
	SQSMockEvent = events.SQSEvent{
		Records: []events.SQSMessage{
			{
				EventSource: "aws:sqs",
				Body:        Payload,
				MessageId:   "1",
			},
			{
				EventSource: "aws:sqs",
				Body:        Payload,
				MessageId:   "2",
			},
		},
	}
	BatchSQSMockEvent = events.SQSEvent{
		Records: []events.SQSMessage{
			{
				EventSource: "aws:sqs",
				Body:        Payload,
				MessageId:   "1",
			},
			{
				EventSource: "aws:sqs",
				Body:        BadPayload,
				MessageId:   "2",
			},
		},
	}
)
