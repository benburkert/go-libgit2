package libgit2

/*
#cgo pkg-config: --static libgit2
#cgo LDFLAGS: -lgit2
#include "libgit2.h"
*/
import "C"

import "unsafe"

// Init initializes libgit2 global state. Init must be called before any other
// libgit2 function in order to set up global state and threading.
func Init() {
	C.git_libgit2_init()
}

// Shutdown cleans up libgit2 global state.
func Shutdown() {
	C.git_libgit2_shutdown()
}

func cbool(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}

func ucbool(b bool) C.uint {
	if b {
		return C.uint(1)
	}
	return C.uint(0)
}

func unwrapErr(res C.struct_libgit2_result) error {
	code := C.int(res.code)
	if code >= 0 {
		return nil
	}

	defer C.free(unsafe.Pointer(res.err))

	return &gitError{
		message: C.GoString(res.err.message),
		class:   errorClass(res.err.klass),
		code:    errorCode(code),
	}
}
