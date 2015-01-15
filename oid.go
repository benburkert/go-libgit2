package libgit2

//#include "libgit2.h"
import "C"

type OID struct {
	*gitOID
}

type gitOID struct {
	ptr *C.git_oid
}
