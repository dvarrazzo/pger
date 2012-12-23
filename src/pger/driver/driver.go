package driver

import "C"
import "database/sql"
import "database/sql/driver"
import "errors"

func init() {
	d := new(PgDriver)
	sql.Register("postgresql", d)
}

func (c *PgConn) Begin() (driver.Tx, error) {
	return nil, errors.New("pger: Begin not implemented")
}

func (c *PgConn) Close() error {
	return errors.New("pger: Close not implemented")
}

func (c *PgConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("pger: Prepare not implemented")
}

type PgDriver struct {
}

func (d PgDriver) Open(name string) (driver.Conn, error) {
	return connect(name)
}
