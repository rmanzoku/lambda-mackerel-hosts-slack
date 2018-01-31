package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

// Event from event.json
type Event map[string]interface{}

// HandleRequest is main handler
func HandleRequest(ctx context.Context, event Event) (string, error) {
	return "ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}
