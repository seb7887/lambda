package lambda

import "github.com/aws/aws-lambda-go/lambda"

func StartAPI(svc Service) {
	lambda.Start(New(svc).APIEventHandler)
}

func StartSNS(svc Service) {
	lambda.Start(New(svc).SNSEventHandler)
}

func StartSQS(svc Service) {
	lambda.Start(New(svc).SQSEventHandler)
}

func StartBatchSQS(svc Service) {
	lambda.Start(New(svc).BatchSQSEventHandler)
}

func StartDDB(svc Service) {
	lambda.Start(New(svc).DDBEventHandler)
}

func StartBatchDDB(svc Service) {
	lambda.Start(New(svc).BatchDDBEventHandler)
}
