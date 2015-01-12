package libgit2

/*
#cgo pkg-config: --static libgit2
#cgo LDFLAGS: -lgit2
#include "libgit2.h"
*/
import "C"
import "unsafe"

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
