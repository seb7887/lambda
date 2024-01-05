package testutils

import "github.com/aws/aws-lambda-go/events"

var (
	_payload     = `{"user_uuid":"iao","operation_id":"123"}`
	_badPayload  = `{"user_uuid":"x","operation_id":"123"}`
	SNSMockEvent = events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				EventSource: "aws:sns",
				SNS: events.SNSEntity{
					Message: _payload,
				},
			},
			{
				EventSource: "aws:sns",
				SNS: events.SNSEntity{
					Message: _payload,
				},
			},
		},
	}
	SQSMockEvent = events.SQSEvent{
		Records: []events.SQSMessage{
			{
				EventSource: "aws:sqs",
				Body:        _payload,
				MessageId:   "1",
			},
			{
				EventSource: "aws:sqs",
				Body:        _payload,
				MessageId:   "2",
			},
		},
	}
	BatchSQSMockEvent = events.SQSEvent{
		Records: []events.SQSMessage{
			{
				EventSource: "aws:sqs",
				Body:        _payload,
				MessageId:   "1",
			},
			{
				EventSource: "aws:sqs",
				Body:        _badPayload,
				MessageId:   "2",
			},
		},
	}
)
