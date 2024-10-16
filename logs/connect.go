package logs

import (
	"context"

	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/msg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() (*Conn, error) {
	ctx := context.TODO()
	host := envar.GetStr("", "MONGO_HOST")
	password := envar.GetStr("", "MONGO_PASSWORD")
	dbname := envar.GetStr("data", "MONGO_DB")

	if host == "" {
		return nil, Alertf(msg.ERR_ENV_REQUIRED, "MONGO_HOST")
	}

	if password == "" {
		return nil, Alertf(msg.ERR_ENV_REQUIRED, "MONGO_PASSWORD")
	}

	clientOptions := options.Client().ApplyURI(host).SetAuth(options.Credential{Password: password}).SetDirect(true)
	client, err := mongo.Connect(ctx, clientOptions)
	db := client.Database(dbname)
	if err != nil {
		Errorf("Error connecting to MongoDB: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		Errorf("Error connecting to MongoDB: %v", err)
	}

	return &Conn{
		ctx:    ctx,
		host:   host,
		dbname: dbname,
		db:     db,
	}, nil
}
