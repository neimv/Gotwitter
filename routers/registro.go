package routers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/neimv/Gotwitter/db"
	"github.com/neimv/Gotwitter/models"
)

func Registro(ctx context.Context) models.RespApi {
	var t models.Usuario
	var r models.RespApi
	r.Status = 400

	fmt.Println("entre a registro")

	body := ctx.Value(models.Key("body")).(string)

	err := json.Unmarshal([]byte(body), &t)

	if err != nil {
		r.Message = err.Error()
		fmt.Println(r.Message)
		return r
	}

	if len(t.Email) == 0 {
		r.Message = "Debe especificar el Email"
		fmt.Println(r.Message)
	}
	if len(t.Password) < 6 {
		r.Message = "Debe especificar una contrasenia de al menos 6 caracteres"
		fmt.Println(r.Message)
	}

	_, encontrado, _ := db.ChequeoYaExisteUsuario(t.Email)

	if encontrado {
		r.Message = "Ya eiste un usuario registrado con ese email"
		fmt.Println(r.Message)
		return r
	}

	_, status, err := db.InsertoRegistro(t)

	if err != nil {
		r.Message = "Ocurrio un error al intentar realizar el registro del usuario" + err.Error()
		fmt.Println(r.Message)
		return r
	}

	if !status {
		r.Message = "No se ha logrado insertar el registro del usuario"
		fmt.Println(r.Message)
		return r
	}

	r.Status = 200
	r.Message = "registro OK"
	fmt.Println(r.Message)

	return r
}
