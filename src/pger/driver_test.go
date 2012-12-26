package pger

import "os"
import "testing"
import "database/sql"

import _ "pger/driver"

var conninfo string

func init() {
	conninfo = os.Getenv("PGER_TESTDSN")
}

func TestConnect(t *testing.T) {
	cnn, err := sql.Open("postgresql", conninfo)
	if err != nil {
		t.Errorf("connection failed: %v", err)
		return
	}

	rows, err := cnn.Query("select 42::int8")
	if err != nil {
		t.Errorf("query failed: %v", err)
		return
	}

	for rows.Next() {
		var n int
		err = rows.Scan(&n)
		if err != nil {
			t.Errorf("scan failed: %v", err)
			return
		}
		if n != 42 {
			t.Errorf("scan failed: n=%v", n)
			return
		}
	}
}
