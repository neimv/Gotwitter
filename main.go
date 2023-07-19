package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"

	"github.com/neimv/Gotwitter/awsgo"
)

func main() {
	lambda.Start(EjecutoLambda)
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAws()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error en las variables de entorno (SecretName, BucketName, UrlPrefix)",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	}

	return res, nil
}

func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")

	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("BucketName")

	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("UrlPrefix")

	if !traeParametro {
		return traeParametro
	}

	return traeParametro
}
