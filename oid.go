package libgit2

//#include "libgit2.h"
import "C"

type OID struct {
	*gitOID
}

func (o OID) String() string {
	return gitOIDTostr(o.gitOID)
}

func (o OID) isZero() bool {
	return gitOIDIszero(o.gitOID)
}

type gitOID struct {
	ptr *C.git_oid
}

func gitOIDIszero(oid *gitOID) bool {
	return int(C.git_oid_iszero(oid.ptr)) == 1
}

func gitOIDTostr(oid *gitOID) string {
	nbuf := C.GIT_OID_HEXSZ + 1
	buf := make([]C.char, nbuf)
	return C.GoString(C.git_oid_tostr(&buf[0], C.size_t(nbuf), oid.ptr))
}
