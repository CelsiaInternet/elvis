package logs

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var conn *Conn

type Conn struct {
	ctx    context.Context
	host   string
	dbname int
	db     *mongo.Client
}

func Load() (*Conn, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Close() error {
	if conn.db == nil {
		return nil
	}

	return conn.db.Disconnect(conn.ctx)
}
