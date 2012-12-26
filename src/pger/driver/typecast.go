package driver

import "C"

import "strconv"

var castmap = map[int]func(*C.char, C.int, *PgRows) (interface{}, error){
	20: castInt,
	// TODO: default
	// TODO: all supported type
	// TODO: extendibility
}

func typecast(oid int, data *C.char, length C.int, rows *PgRows) (interface{}, error) {
	f, ok := castmap[oid]
	if !ok {
		f = castBytes
	}
	return f(data, length, rows)
}

func castBytes(data *C.char, length C.int, rows *PgRows) (interface{}, error) {
	s := C.GoStringN(data, length)
	return []byte(s), nil
}

func castInt(data *C.char, length C.int, rows *PgRows) (interface{}, error) {
	s := C.GoStringN(data, length)
	return strconv.ParseInt(s, 10, 64)
}
