package jdb

/**
* Load
* @return *Conn, error
**/
func Load() (*DB, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}

	if !conn.UseCore {
		return conn, nil
	}

	err = InitCore(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
