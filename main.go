package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"

	"github.com/neimv/Gotwitter/awsgo"
	"github.com/neimv/Gotwitter/db"
	"github.com/neimv/Gotwitter/handlers"
	"github.com/neimv/Gotwitter/models"
	"github.com/neimv/Gotwitter/secretmanager"
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

	secretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	fmt.Println("Secretro obtenido")
	fmt.Println(err)

	if err != nil {
		fmt.Println(err)
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error en la lectura de Secrets" + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	}

	path := strings.Replace(request.PathParameters["twittergo"], os.Getenv("UrlPrefix"), "", -1)

	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), secretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), secretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), secretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), secretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), secretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	// chequeo de conexion a la base de datos
	err = db.ConectarDB(awsgo.Ctx)

	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error en la Conexion de base de datos" + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	respAPI := handlers.Manejadores(awsgo.Ctx, request)

	if respAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       respAPI.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	} else {
		return respAPI.CustomResp, nil
	}

	// return res, nil
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
