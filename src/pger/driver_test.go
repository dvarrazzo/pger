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

	rows, err := cnn.Query(
		"select $1::int8 as foo, $2::int8 as bar", 42, nil)
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

	nrows := 0
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
		nrows++
	}

	if nrows != 1 {
		t.Fatal("bad rows number:", nrows)
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
		create table test_exec (
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

	res, err = cnn.Exec("insert into test_exec (data) values ($1)", "hello")
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

func TestEnc(t *testing.T) {
	cnn, err := sql.Open("postgresql", conninfo)
	if err != nil {
		t.Fatal("connection failed:", err)
	}
	defer cnn.Close()

	if _, err := cnn.Exec("drop table if exists test_enc"); err != nil {
		t.Fatal("drop table failed:", err)
	}
	if _, err := cnn.Exec(`
		create table test_enc (
			id serial primary key, data text)`); err != nil {
		t.Fatal("create table failed:", err)
	}
	// right arg encoding
	if _, err := cnn.Exec(
		"insert into test_enc (data) values ($1)", "€"); err != nil {
		t.Fatal("insert 1 failed:", err)
	}
	// right query encoding
	if _, err := cnn.Exec(
		"insert into test_enc (data) values ('☃')"); err != nil {
		t.Fatal("insert 2 failed:", err)
	}

	for id, v := range map[int]string{1: "€", 2: "☃"} {
		row := cnn.QueryRow("select data from test_enc where id = $1", id)
		var s string
		if err := row.Scan(&s); err != nil {
			t.Fatal("Scan failed:", err)
		}
		if s != v {
			t.Fatal("bad queried value for id", id, ":", s, "instead of", v)
		}
	}
}

func TestTx(t *testing.T) {
	cnn, err := sql.Open("postgresql", conninfo)
	if err != nil {
		t.Fatal("connection failed:", err)
	}
	defer cnn.Close()

	if _, err := cnn.Exec("drop table if exists test_tx"); err != nil {
		t.Fatal("drop table failed:", err)
	}
	if _, err := cnn.Exec(`
		create table if not exists test_tx (
			id serial primary key, data text)`); err != nil {
		t.Fatal("create table failed:", err)
	}

	tx, err := cnn.Begin()
	if err != nil {
		t.Fatal("begin 1 failed:", err)
	}
	if _, err := tx.Exec(
		"insert into test_tx (data) values ($1)", "foo"); err != nil {
		t.Fatal("insert 1 failed:", err)
	}
	if err = tx.Rollback(); err != nil {
		t.Fatal("rollback failed:", err)
	}

	tx, err = cnn.Begin()
	if err != nil {
		t.Fatal("begin 2 failed:", err)
	}
	if _, err := tx.Exec(
		"insert into test_tx (data) values ($1)", "bar"); err != nil {
		t.Fatal("insert 2 failed:", err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal("commit failed:", err)
	}

	row := cnn.QueryRow("select count(*) from test_tx")
	var n int
	if err := row.Scan(&n); err != nil {
		t.Fatal("Scan 1 failed:", err)
	}
	if n != 1 {
		t.Fatal("expected 1 record; found:", n)
	}

	row = cnn.QueryRow("select data from test_tx where id = 2")
	var s string
	if err := row.Scan(&s); err != nil {
		t.Fatal("Scan 2 failed:", err)
	}
	if s != "bar" {
		t.Fatal("expected 'bar' record; found:", s)
	}
}
