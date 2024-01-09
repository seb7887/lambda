package lambda

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/tidwall/gjson"
)

/*
*
aws:dynamodb Records.0.eventSource
aws:sqs      Records.0.eventSource
aws:s3       Records.0.eventSource
aws:sns      Records.0.EventSource
aws:cw       source
aws:apigw    default
*/
const (
	_ddbEvent = "aws:dynamodb"
	_sqsEvent = "aws:sqs"
	_s3Event  = "aws:s3"
	_snsEvent = "aws:sns"
	_cwEvent  = "aws:cw"
	_apiEvent = "aws:apigw"
)

var eventKeys = []string{
	"Records.0.eventSource",
	"Records.0.EventSource",
}

func eventSource(raw json.RawMessage) string {
	jsonStr := string(raw)

	for _, key := range eventKeys {
		if source := gjson.Get(jsonStr, key); source.Exists() {
			return source.String()
		}
	}

	if source := gjson.Get(jsonStr, "source"); source.Exists() {
		return _cwEvent
	}

	return _apiEvent
}

func identifyEventType(raw json.RawMessage) (any, error) {
	switch eventSource(raw) {
	case _sqsEvent:
		var evt events.SQSEvent
		if err := json.Unmarshal(raw, &evt); err != nil {
			return nil, err
		}
		return &evt, nil
	case _snsEvent:
		var evt events.SNSEvent
		if err := json.Unmarshal(raw, &evt); err != nil {
			return nil, err
		}
		return &evt, nil
	case _ddbEvent:
		var evt events.DynamoDBEvent
		if err := json.Unmarshal(raw, &evt); err != nil {
			return nil, err
		}
		return &evt, nil
	case _cwEvent:
		var evt events.CloudWatchEvent
		if err := json.Unmarshal(raw, &evt); err != nil {
			return nil, err
		}
		return &evt, nil
	case _s3Event:
		var evt events.S3Event
		if err := json.Unmarshal(raw, &evt); err != nil {
			return nil, err
		}
		return &evt, nil
	default:
		return nil, errors.New("invalid event")
	}
}
