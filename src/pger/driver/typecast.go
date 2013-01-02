package driver

import "C"

import "strconv"

type castfunc func(*C.char, C.int, *PgRows) (interface{}, error)

var castmap = map[int]castfunc{
	20: castInt,
	// TODO: default
	// TODO: all supported type
	// TODO: extendibility
}

func typecast(oid int, data *C.char, length C.int, rows *PgRows) (interface{}, error) {
	f, ok := castmap[oid]
	if !ok {
		f = castString
	}
	return f(data, length, rows)
}

func castString(data *C.char, length C.int, rows *PgRows) (interface{}, error) {
	s := C.GoStringN(data, length)
	return s, nil
}

func castInt(data *C.char, length C.int, rows *PgRows) (interface{}, error) {
	s := C.GoStringN(data, length)
	return strconv.ParseInt(s, 10, 64)
}
