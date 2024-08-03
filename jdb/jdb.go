package jdb

/**
* Ths jdb package makes it easy to create an array of database connections
* initially to posrtgresql databases.
*	Provide a connection function, validate the existence of elements such as databases, schemas, tables, colums, index, series and users and
* it is possible to create them if they do not exist.
* Also, have a execute to sql sentences to retuns json and json array,
* that valid you result return records and how many records are returned.
**/

var (
	conn *Conn
)

func Load() (*Conn, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = connect()
	if err != nil {
		return nil, err
	}

	err = InitCore(conn.Db)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Close() error {
	err := conn.Db.Close()
	if err != nil {
		return err
	}

	return nil
}
