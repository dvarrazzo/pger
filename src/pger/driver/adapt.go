package driver

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

func AdaptValues(args []driver.Value, c *PgConn) ([]*string, error) {
	rv := make([]*string, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == nil {
			rv[i] = nil
		} else {
			var err error
			if rv[i], err = Adapt(args[i], c); err != nil {
				return nil, err
			}
		}
	}
	return rv, nil
}

func Adapt(v interface{}, c *PgConn) (*string, error) {
	switch t := v.(type) {
	case int64:
		return AdaptInt64(t, c)
	case string:
		return AdaptString(t, c)
	default:
		return nil, fmt.Errorf("can't adapt type %T: %v", t, t)
	}
	return nil, fmt.Errorf("wat")
}

func AdaptInt64(v int64, c *PgConn) (*string, error) {
	s := strconv.FormatInt(v, 10)
	return &s, nil
}

func AdaptString(s string, c *PgConn) (*string, error) {
	return &s, nil
}
