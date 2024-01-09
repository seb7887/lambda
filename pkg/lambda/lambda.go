package lambda

import "github.com/aws/aws-lambda-go/lambda"

func StartAPI(svc HandlerService) {
	lambda.Start(New(svc).APIEventHandler)
}

func StartSNS(svc HandlerService) {
	lambda.Start(New(svc).SNSEventHandler)
}

func StartSQS(svc HandlerService) {
	lambda.Start(New(svc).SQSEventHandler)
}

func StartBatchSQS(svc HandlerService) {
	lambda.Start(New(svc).BatchSQSEventHandler)
}

func StartDDB(svc HandlerService) {
	lambda.Start(New(svc).DDBEventHandler)
}

func StartBatchDDB(svc HandlerService) {
	lambda.Start(New(svc).BatchDDBEventHandler)
}

func StartMulti(svc HandlerService) {
	lambda.Start(New(svc).MultiEventHandler)
}

func Start[I, O any](fn HandlerFn[I, O]) {
	lambda.Start(NewHandler(fn, LoggerMiddleware[I, O]).EventHandler)
}
