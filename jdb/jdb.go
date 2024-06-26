package jdb

/**
* Ths jdb package makes it easy to create an array of database connections
* initially to posrtgresql databases.
*	Provide a connection function, validate the existence of elements such as databases, schemas, tables, colums, index, series and users and
* it is possible to create them if they do not exist.
* Also, have a execute to sql sentences to retuns json and json array,
* that valid you result return records and how many records are returned.
**/

type Conn struct {
	Db []*Db
}

var conn *Conn

func Load() (*Conn, error) {
	if conn != nil {
		return conn, nil
	}

	db, err := connect()
	if err != nil {
		return nil, err
	}

	if conn == nil {
		conn = &Conn{
			Db: []*Db{},
		}
	}

	idx := len(conn.Db)
	db.Index = idx

	conn.Db = append(conn.Db, db)

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

func DB(idx int) *Db {
	return conn.Db[idx]
}

func DBClose(idx int) error {
	return conn.Db[idx].Close()
}
