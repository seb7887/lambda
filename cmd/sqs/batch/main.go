package main

import (
	"seb7887/lambda/internal/service"
	"seb7887/lambda/pkg/lambda"
)

func main() {
	lambda.StartBatchSQS(service.New())
}
