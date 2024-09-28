package jdb

const (
	CommandDefine = "DEFINE"
	CommandInsert = "INSERT"
	CommandUpdate = "UPDATE"
	CommandDelete = "DELETE"
)

/**
* Load
* @return *Conn, error
**/
func Load() (*DB, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}

	err = InitCore(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
