package driver

/*
// TODO: configure include dir
#cgo CFLAGS: -I/usr/include/postgresql
#cgo LDFLAGS: -lpq
#include <libpq-fe.h>
*/
import "C"

import "fmt"

type PgConn struct {
	conninfo string
	conn     *C.PGconn
}

func connect(conninfo string) (*PgConn, error) {
	cs := C.CString(conninfo)
	conn := C.PQconnectdb(cs)
	if C.PQstatus(conn) != C.CONNECTION_OK {
		cerr := C.PQerrorMessage(conn)
		err := fmt.Errorf("connection failed: %s", C.GoString(cerr))
		C.PQfinish(conn)
		return nil, err
	}

	return &PgConn{conn: conn, conninfo: conninfo}, nil
}
