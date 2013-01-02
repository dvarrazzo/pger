package driver

/*
// TODO: configure include dir
#cgo CFLAGS: -I/usr/include/postgresql
#cgo LDFLAGS: -lpq
#include <libpq-fe.h>
#include <stdlib.h>
*/
import "C"

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

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
	tuples  string
	lastoid uint
}

type PgStmt struct {
	name   string
	result *C.PGresult
	conn   *PgConn
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

	stmt := PgStmt{conn: conn, result: res, name: ""}
	return &stmt, nil
}

type PgRows struct {
	result *C.PGresult
	cur    int
}

func query(stmt PgStmt, values []driver.Value) (*PgRows, error) {
	res, err := exec_prepared(stmt, values)
	if err != nil {
		return &PgRows{}, err
	}

	rows := PgRows{result: res}
	return &rows, nil
}

func exec(stmt PgStmt, values []driver.Value) (*PgResult, error) {
	res, err := exec_prepared(stmt, values)
	if err != nil {
		return &PgResult{}, err
	}
	defer pqclear(res)

	tuples := C.GoString(C.PQcmdTuples(res))
	lastoid := uint(C.PQoidValue(res))

	rv := PgResult{tuples: tuples, lastoid: lastoid}
	return &rv, nil
}

func exec_prepared(stmt PgStmt, values []driver.Value) (*C.PGresult, error) {
	conn := stmt.conn
	params, err := AdaptValues(values, conn)
	if err != nil {
		return nil, err
	}

	name := C.CString(stmt.name)
	defer C.free(unsafe.Pointer(name))

	cparams := charpp(params)
	defer charppFree(cparams, len(params))

	// TODO: binary params
	// TODO: binary results
	res := C.PQexecPrepared(conn.conn, name, C.int(len(params)),
		cparams, nil, nil, C.int(0))
	if err := errorFromPGresult(res); err != nil {
		pqclear(res)
		return nil, err
	}

	return res, nil
}

func columns(r *PgRows) []string {
	nfields := int(C.PQnfields(r.result))
	rv := make([]string, nfields)
	for i := 0; i < nfields; i++ {
		rv[i] = C.GoString(C.PQfname(r.result, C.int(i)))
	}
	return rv
}

func next(r *PgRows, dest []driver.Value) error {
	if r.cur >= int(C.PQntuples(r.result)) {
		return io.EOF
	}

	// TODO: cache oids
	// TODO: extendibility
	for i := 0; i < len(dest); i++ {
		if 0 == int(C.PQgetisnull(r.result, C.int(r.cur), C.int(i))) {
			oid := int(C.PQftype(r.result, C.int(i)))
			data := C.PQgetvalue(r.result, C.int(r.cur), C.int(i))
			length := C.PQgetlength(r.result, C.int(r.cur), C.int(i))
			var err error
			dest[i], err = typecast(oid, data, length, r)
			if err != nil {
				return err
			}
		} else {
			dest[i] = nil
		}
	}
	r.cur++
	return nil
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
