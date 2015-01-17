package libgit2

//#include "libgit2.h"
import "C"

import (
	"runtime"
	"sync"
)

// Walker is an in-progress walk through the commits in a repo.
type Walker struct {
	*gitRevwalk

	C <-chan *Commit

	err error

	co *sync.Once
	cc chan struct{}
}

func newWalker(config *walkerConfig) (*Walker, error) {
	r, err := gitRevwalkNew(config.repo.gitRepository)
	if err != nil {
		return nil, err
	}

	if config.startRef == "" {
		if err = r.pushHead(); err != nil {
			return nil, err
		}
	}

	c := make(chan *Commit, config.bufSize)
	w := &Walker{
		gitRevwalk: r,
		C:          c,
		co:         &sync.Once{},
		cc:         make(chan struct{}),
	}

	go w.run(config.repo, c)
	return w, nil
}

// Cancel aborts an in-progress walk and drains the commit channel C.
func (w *Walker) Cancel() {
	w.co.Do(w.cancel)
}

// Err returns error encountered while walking commits.
func (w *Walker) Err() error {
	return w.err
}

func (w *Walker) cancel() {
	close(w.cc)
	for range w.C {
	}
}

func (w *Walker) next(repo Repository) (*Commit, error) {
	o, err := w.gitRevwalk.next()
	if err != nil {
		return nil, err
	}

	oid := OID{o}
	if oid.isZero() {
		return nil, nil
	}
	return lookupCommit(repo, oid)
}

func (w *Walker) run(repo Repository, c chan<- *Commit) {
	defer close(c)

	for {
		commit, err := w.next(repo)
		if err != nil {
			w.err = err
			return
		}
		if commit == nil {
			return
		}

		select {
		case c <- commit:
		case <-w.cc:
			return
		}
	}
}

type gitRevwalk struct {
	ptr *C.git_revwalk
}

func (r *gitRevwalk) init() {
	runtime.SetFinalizer(r, (*gitRevwalk).free)
}

func (r *gitRevwalk) next() (*gitOID, error) {
	oid := &gitOID{ptr: &C.git_oid{}}
	return oid, unwrapErr(C.libgit2_revwalk_next(oid.ptr, r.ptr))
}

func (r *gitRevwalk) free() {
	runtime.SetFinalizer(r, nil)
	C.git_revwalk_free(r.ptr)
}

func (r *gitRevwalk) pushHead() error {
	return unwrapErr(C.libgit2_revwalk_push_head(r.ptr))
}

func gitRevwalkNew(repo *gitRepository) (*gitRevwalk, error) {
	var ptr *C.git_revwalk
	if err := unwrapErr(C.libgit2_revwalk_new(&ptr, repo.ptr)); err != nil {
		return nil, err
	}

	r := &gitRevwalk{ptr}
	r.init()
	return r, nil
}
