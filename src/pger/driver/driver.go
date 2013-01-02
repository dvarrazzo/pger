package driver

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"
)

func init() {
	d := new(PgDriver)
	sql.Register("postgresql", d)
}

// Implementation of the Conn interface

func (c *PgConn) Begin() (driver.Tx, error) {
	return nil, errors.New("pger: Begin not implemented")
}

func (c *PgConn) Close() error {
	return close_conn(c)
}

func (c *PgConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := prepare(c, query)
	return *stmt, err
}

// Implementation of the Stmt interface

func (s PgStmt) Close() error {
	pqclear(s.result)
	return nil
}

func (s PgStmt) NumInput() int {
	return -1
}

func (s PgStmt) Exec(args []driver.Value) (driver.Result, error) {
	return exec(s, args)
}

func (s PgStmt) Query(args []driver.Value) (driver.Rows, error) {
	return query(s, args)
}

// Implementation of the Result interface

func (r PgResult) LastInsertId() (int64, error) {
	if r.lastoid == 0 {
		return 0, errors.New("command returned no last insert id")
	}
	return int64(r.lastoid), nil
}

func (r PgResult) RowsAffected() (int64, error) {
	if r.tuples == "" {
		return 0, errors.New("command returned no rows affected number")
	}
	return strconv.ParseInt(r.tuples, 10, 64)
}

// Implementation of the Rows interface

func (r *PgRows) Columns() []string {
	return columns(r)
}

func (r *PgRows) Close() error {
	pqclear(r.result)
	return nil
}

func (r *PgRows) Next(dest []driver.Value) error {
	return next(r, dest)
}

// Implementation of the Driver interface

type PgDriver struct {
}

func (d PgDriver) Open(name string) (driver.Conn, error) {
	return connect(name)
}
