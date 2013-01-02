package pger

import (
	"database/sql"
	"os"
	"testing"
)

import _ "pger/driver"

var conninfo string

func init() {
	conninfo = os.Getenv("PGER_TESTDSN")
}

func TestQuery(t *testing.T) {
	cnn, err := sql.Open("postgresql", conninfo)
	if err != nil {
		t.Fatal("connection failed:", err)
	}
	defer cnn.Close()

	rows, err := cnn.Query("select 42::int8 as foo, null::int8 as bar")
	if err != nil {
		t.Fatal("query failed:", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		t.Fatal("rows.Columns() failed:", err)
	}
	if len(cols) != 2 {
		t.Fatal("bad number of columns:", len(cols))
	}
	if cols[0] != "foo" || cols[1] != "bar" {
		t.Fatal("bad column names:", cols)
	}

	for rows.Next() {
		var n1 int
		var n2 sql.NullInt64
		err = rows.Scan(&n1, &n2)
		if err != nil {
			t.Fatal("scan failed:", err)
		}
		if n1 != 42 {
			t.Fatal("scan failed: n1 =", n1)
		}
		if n2.Valid {
			t.Fatal("scan failed: n2 =", n2.Int64)
		}
	}
}

func TestExec(t *testing.T) {
	cnn, err := sql.Open("postgresql", conninfo)
	if err != nil {
		t.Fatal("connection failed:", err)
	}
	defer cnn.Close()

	res, err := cnn.Exec("drop table if exists test_exec")
	if err != nil {
		t.Fatal("drop table failed:", err)
	}

	res, err = cnn.Exec(`
		create table if not exists test_exec (
			id serial primary key, data text) with oids`)
	if err != nil {
		t.Fatal("create table failed:", err)
	}
	n, rerr := res.LastInsertId()
	if rerr == nil {
		t.Fatal("no error after LastInsertId on create table")
	}
	n, rerr = res.RowsAffected()
	if rerr == nil {
		t.Fatal("no error after RowsAffected on create table")
	}

	res, err = cnn.Exec("insert into test_exec (data) values ('hello')")
	if err != nil {
		t.Fatal("exec failed:", err)
	}
	n, rerr = res.LastInsertId()
	if rerr != nil {
		t.Fatal("error after LastInsertId on insert:", err)
	}
	if n == 0 {
		t.Fatal("invalid oid returned after insert")
	}
	n, rerr = res.RowsAffected()
	if rerr != nil {
		t.Fatal("error after RowsAffected on insert")
	}
	if n != 1 {
		t.Fatal("bad rows affected after single insert:", n)
	}
}
