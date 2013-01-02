package driver

/*
#include <stdlib.h>		// malloc
#include <string.h>		// memset

static char **
charpp_make(int size)
{
	char **buf;
	buf = (char **)malloc(size * sizeof(char *));
	memset((void *)buf, 0, size * sizeof(char *));
	return buf;
}

static void
charpp_free(char **args, int size) {
	int i;
	for (i = 0; i < size; i++) {
		if (args[i])
			free(args[i]);
	}
	free(args);
}

static void
charpp_set(char **args, int i, char *val)
{
	if (args[i])
		free(args[i]);
	args[i] = val;
}
*/
import "C"

func charpp(args [][]byte) **C.char {
	buf := C.charpp_make(C.int(len(args)))
	return buf
}

func charppFree(args **C.char, n int) {
	C.charpp_free(args, C.int(n))
}

func charppSet(args **C.char, i int, val []byte) {
	// TODO: looks like adapt should return string, not []byte
	C.charpp_set(args, C.int(i), C.CString(string(val)))
}
