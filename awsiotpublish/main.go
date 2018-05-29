package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"net/url"
)

type IotMessage struct {
	Message string `json:"message"`
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	logObject(event)

	region := "eu-west-1"
	awsConfig := aws.Config{
		Region: aws.String(region),
	}
	awsConfig.WithLogLevel(aws.LogDebugWithRequestErrors)
	session := session.New(&awsConfig)

	topic := event.PathParameters["topic"]
	messageUnescaped, err := url.QueryUnescape(event.PathParameters["message"])
	if err != nil {
		panic(err.Error())
	}
	iotMessage := IotMessage{
		Message: messageUnescaped,
	}
	message, err := json.Marshal(&iotMessage)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(fmt.Sprintf("Sending message '%s' to the topic '%s'", string(message), topic))

	// note: getting the endpoint address can be skipped and address hard coded for slightly better performance
	iotSvc := iot.New(session)
	res, err := iotSvc.DescribeEndpoint(&iot.DescribeEndpointInput{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(fmt.Sprintf("Endpoint address: %s", *res.EndpointAddress))

	svc := iotdataplane.New(session, &aws.Config{Endpoint: res.EndpointAddress})

	publishInput := iotdataplane.PublishInput{
		Topic:   aws.String(topic),
		Payload: message,
		Qos:     aws.Int64(1),
	}

	publishOutput, err := svc.Publish(&publishInput)

	var responseMessage []byte
	var responseCode int
	if err != nil {
		fmt.Println(fmt.Sprintf(err.Error()))
		responseMessage, _ = json.Marshal(err.Error())
		responseCode = 500
	} else {
		responseMessage, _ = json.Marshal(publishOutput)
		responseCode = 200
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseMessage),
		StatusCode: responseCode,
	}, nil
}

func logObject(object interface{}) {
	result, err := json.Marshal(object)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(fmt.Sprintf("%s", result))
}

func main() {
	lambda.Start(Handler)
}
