package libgit2

//#include "libgit2.h"
import "C"
import "runtime"

type Reference struct {
	*gitReference
}

func (r *Reference) target() *OID {
	return &OID{gitReferenceTarget(r.gitReference)}
}

type gitReference struct {
	ptr *C.git_reference
}

func (r *gitReference) init() {
	runtime.SetFinalizer(r, (*gitReference).free)
}

func (r *gitReference) free() {
	runtime.SetFinalizer(r, nil)
	C.git_reference_free(r.ptr)
}

func gitReferenceTarget(ref *gitReference) *gitOID {
	return &gitOID{ptr: C.git_reference_target(ref.ptr)}
}
