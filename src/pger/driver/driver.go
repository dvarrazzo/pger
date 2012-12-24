package driver

import "database/sql"
import "database/sql/driver"
import "errors"

func init() {
	d := new(PgDriver)
	sql.Register("postgresql", d)
}

// Implementation of the Conn interface

func (c *PgConn) Begin() (driver.Tx, error) {
	return nil, errors.New("pger: Begin not implemented")
}

func (c *PgConn) Close() error {
	return errors.New("pger: Close not implemented")
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
	return PgResult{}, errors.New("pger: Exec not implemented")
}

func (s PgStmt) Query(args []driver.Value) (driver.Rows, error) {
	return PgRows{}, errors.New("pger: Query not implemented")
}

// Implementation of the Result interface

func (r PgResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r PgResult) RowsAffected() (int64, error) {
	return 1, nil
}

// Implementation of the Rows interface

func (r PgRows) Columns() []string {
	// TODO
	return []string{}
}

func (r PgRows) Close() error {
	// TODO
	return nil
}

func (r PgRows) Next(dest []driver.Value) error {
	// TODO
	return nil
}

// Implementation of the Driver interface

type PgDriver struct {
}

func (d PgDriver) Open(name string) (driver.Conn, error) {
	return connect(name)
}
