package logs

import (
	"context"

	"github.com/cgalvisleon/elvis/envar"
	"github.com/cgalvisleon/elvis/msg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() (*Conn, error) {
	ctx := context.TODO()
	host := envar.EnvarStr("", "MONGO_HOST")
	password := envar.EnvarStr("", "MONGO_PASSWORD")
	dbname := envar.EnvarInt(0, "MONGO_DB")

	if host == "" {
		return nil, Alertf(msg.ERR_ENV_REQUIRED, "MONGO_HOST")
	}

	if password == "" {
		return nil, Alertf(msg.ERR_ENV_REQUIRED, "MONGO_PASSWORD")
	}

	client := options.Client().ApplyURI(host).SetAuth(options.Credential{Password: password}).SetDirect(true)
	clients, err := mongo.Connect(ctx, client)
	if err != nil {
		Errorf("Error connecting to MongoDB: %v", err)
	}
	err = clients.Ping(ctx, nil)
	if err != nil {
		Errorf("Error connecting to MongoDB: %v", err)
	}

	return &Conn{
		ctx:    ctx,
		host:   host,
		dbname: dbname,
		db:     clients,
	}, nil
}
