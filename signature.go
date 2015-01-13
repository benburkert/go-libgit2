package libgit2

//#include "libgit2.h"
import "C"

import (
	"runtime"
	"time"
)

type Signature struct {
	*gitSignature
}

func defaultSignature(repo Repository) (*Signature, error) {
	sig, err := gitSignatureDefault(repo.gitRepository)
	if err != nil {
		return nil, err
	}
	return &Signature{sig}, nil
}

type gitSignature struct {
	ptr *C.git_signature

	// Name is the full name of the author.
	Name string
	// Email is the email of the author.
	Email string
	// When is the time time an action happened.
	When time.Time
}

func (s *gitSignature) init() {
	// git stores minutes, go wants seconds
	loc := time.FixedZone("", int(s.ptr.when.offset)*60)

	s.Name = C.GoString(s.ptr.name)
	s.Email = C.GoString(s.ptr.email)
	s.When = time.Unix(int64(s.ptr.when.time), 0).In(loc)

	runtime.SetFinalizer(s, (*gitSignature).free)
}

func (s *gitSignature) free() {
	runtime.SetFinalizer(s, nil)
	C.git_signature_free(s.ptr)
}

func gitSignatureDefault(repo *gitRepository) (*gitSignature, error) {
	s := new(gitSignature)

	err := unwrapErr(C.libgit2_signature_default(&s.ptr, repo.ptr))
	if err != nil {
		return nil, err
	}
	s.init()
	return s, nil
}
