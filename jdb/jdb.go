package jdb

var conn *DB

/**
* Load
* @return *Conn, error
**/
func Load() (*DB, error) {
	if conn != nil {
		return conn, nil
	}

	var err error
	conn, err = connect()
	if err != nil {
		return nil, err
	}

	err = InitCore(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

/**
* Close
* @return error
**/
func Close() error {
	if conn == nil {
		return nil
	}

	if conn.db != nil {
		err := conn.db.Close()
		if err != nil {
			return err
		}
	}

	if conn.dm != nil {
		err := conn.dm.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
