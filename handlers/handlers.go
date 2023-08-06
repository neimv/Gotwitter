package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/neimv/Gotwitter/jwt"
	"github.com/neimv/Gotwitter/models"
	"github.com/neimv/Gotwitter/routers"
)

func Manejadores(ctx context.Context, request events.APIGatewayProxyRequest) models.RespApi {
	fmt.Println("Voy a procesar " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var r models.RespApi

	isOk, statusCode, msg, claim := validoAuthorization(ctx, request)
	fmt.Println(claim)

	if !isOk {
		r.Status = statusCode
		r.Message = msg

		return r
	}

	r.Status = 400

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		case "registro":
			return routers.Registro(ctx)
		}
		//
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {

		}
		//
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {

		}
		//
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {

		}
		//
	}

	r.Message = "Method invalid"

	return r
}

func validoAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	path := ctx.Value(models.Key("path")).(string)

	if path == "registro" || path == "login" || path == "obtenerAvatar" || path == "obtenerBanner" {
		return true, 200, "", models.Claim{}
	}

	token := request.Headers["Authorization"]

	if len(token) == 0 {
		return false, 401, "token requerido", models.Claim{}
	}

	claim, todoOk, msg, err := jwt.ProcesoToken(token, ctx.Value(models.Key("jwtSign")).(string))

	if !todoOk {
		if err != nil {
			fmt.Println("error en el token " + err.Error())

			return false, 401, err.Error(), models.Claim{}
		} else {
			fmt.Println("Error en el token " + msg)
			return false, 401, msg, models.Claim{}
		}
	}

	fmt.Println("token OK")

	return true, 200, msg, *claim
}
