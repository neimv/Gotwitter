package db

import (
	"context"
	"fmt"

	"github.com/neimv/Gotwitter/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoCN *mongo.Client
var DatabaseName string

func ConectarDB(ctx context.Context) error {
	user := ctx.Value(models.Key("user")).(string)
	password := ctx.Value(models.Key("password")).(string)
	host := ctx.Value(models.Key("host")).(string)
	connStr := fmt.Sprintf("mongdb+srv://%s:%s@%s/?retryWrites=true", user, password, host)

	var clientOptions = options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		fmt.Println(err.Error())

		return err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		fmt.Println(err.Error())

		return err
	}

	fmt.Println("Conexion Existosa con la DB")
	MongoCN = client
	DatabaseName = ctx.Value(models.Key("database")).(string)

	return nil
}

func BaseConectada() bool {
	err := MongoCN.Ping(context.TODO(), nil)

	return err == nil
}
