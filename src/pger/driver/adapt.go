package driver

import "database/sql/driver"
import "fmt"
import "strconv"

func AdaptValues(args []driver.Value, c *PgConn) ([][]byte, error) {
	rv := make([][]byte, len(args))
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

func Adapt(v interface{}, c *PgConn) ([]byte, error) {
	switch t := v.(type) {
	case int64:
		return AdaptInt64(t, c)
	default:
		return nil, fmt.Errorf("can't adapt type %T: %v", t, t)
	}
	return nil, fmt.Errorf("wat")
}

func AdaptInt64(v int64, c *PgConn) ([]byte, error) {
	s := strconv.FormatInt(v, 10)
	return []byte(s), nil
}
