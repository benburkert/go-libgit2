package libgit2

//#include "libgit2.h"
import "C"
import "unsafe"

func gitMessagePrettify(msg string, strip bool, char rune) (string, error) {
	buf := &C.git_buf{}
	defer C.git_buf_free(buf)

	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	err := unwrapErr(C.libgit2_message_prettify(buf, cmsg, cbool(strip),
		C.char(char)))
	if err != nil {
		return "", err
	}
	return C.GoString(buf.ptr), nil
}
