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

func charpp(args []*string) **C.char {
	buf := C.charpp_make(C.int(len(args)))
	for i := 0; i < len(args); i++ {
		if args[i] != nil {
			C.charpp_set(buf, C.int(i), C.CString(*args[i]))
		}
	}
	return buf
}

func charppFree(args **C.char, n int) {
	C.charpp_free(args, C.int(n))
}
