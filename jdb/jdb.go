package jdb

import (
	"sync"
)

var (
	conn *Conn
	once sync.Once
)

type Conn struct {
	Db []*Db
}

func Load() (*Conn, error) {
	once.Do(connect)

	return conn, nil
}

func Close() error {
	for _, db := range conn.Db {
		err := db.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func DB(db int) *Db {
	return conn.Db[db]
}

func DBClose(db int) error {
	return conn.Db[db].Close()
}

func Jdb(idx int) *Db {
	return conn.Db[idx]
}
