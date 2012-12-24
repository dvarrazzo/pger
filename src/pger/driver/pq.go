package driver

/*
// TODO: configure include dir
#cgo CFLAGS: -I/usr/include/postgresql
#cgo LDFLAGS: -lpq
#include <libpq-fe.h>
#include <stdlib.h>
*/
import "C"

import "errors"
import "fmt"
import "unsafe"

type PgConn struct {
	conninfo string
	conn     *C.PGconn
}

func connect(conninfo string) (*PgConn, error) {
	cs := C.CString(conninfo)
	defer C.free(unsafe.Pointer(cs))
	conn := C.PQconnectdb(cs)
	if C.PQstatus(conn) != C.CONNECTION_OK {
		cerr := C.PQerrorMessage(conn)
		err := fmt.Errorf("connection failed: %s", C.GoString(cerr))
		C.PQfinish(conn)
		return nil, err
	}

	// one-line error message
	C.PQsetErrorVerbosity(conn, C.PQERRORS_TERSE)

	return &PgConn{conn: conn, conninfo: conninfo}, nil
}

type PgResult struct {
}

type PgStmt struct {
	name string
	result *C.PGresult
}

type PgRows struct {
}

func prepare(conn *PgConn, query string) (*PgStmt, error) {
	name := C.CString("")
	defer C.free(unsafe.Pointer(name))
	cquery := C.CString(query)
	defer C.free(unsafe.Pointer(cquery))
	res := C.PQprepare(conn.conn, name, cquery, 0, nil)

	if err := errorFromPGresult(res); err != nil {
		pqclear(res)
		return &PgStmt{}, err
	}

	stmt := PgStmt{result: res, name: ""}
	return &stmt, nil
}

func errorFromPGresult(res *C.PGresult) error {
	switch C.PQresultStatus(res) {
	case C.PGRES_EMPTY_QUERY:
		return errors.New("empty query")
	case C.PGRES_BAD_RESPONSE:
		return errors.New("bad response")
	case C.PGRES_FATAL_ERROR:
		return errors.New(C.GoString(C.PQresultErrorMessage(res)))
	}
	return nil
}

func pqclear(res *C.PGresult) {
	C.PQclear(res)
}
