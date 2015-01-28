package libgit2

//#include "libgit2.h"
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

type branchType uint

const (
	branchAll    branchType = C.GIT_BRANCH_ALL
	branchLocal  branchType = C.GIT_BRANCH_LOCAL
	branchRemote branchType = C.GIT_BRANCH_REMOTE
)

// Branch represents a local or remote branch in git.
type Branch struct {
	*gitReference
	t branchType
}

func createBranch(config *branchConfig) (*Branch, error) {
	ref, err := gitBranchCreate(config.repo.gitRepository, config.name,
		config.target.gitCommit, config.force, config.sig.gitSignature,
		config.logMessage)
	if err != nil {
		return nil, err
	}
	return &Branch{ref, branchLocal}, nil
}

// Delete an existing branch reference.
func (b Branch) Delete() error {
	return gitBranchDelete(b.gitReference)
}

// Name return the name of the given local or remote branch.
func (b Branch) Name() (string, error) {
	return gitBranchName(b.gitReference)
}

// BranchWalker is an in-progress walk of branches in a repo.
type BranchWalker struct {
	*gitBranchIterator

	C <-chan *Branch

	err error

	co *sync.Once
	cc chan struct{}
}

func newBranchWalker(r Repository, t branchType) (*BranchWalker, error) {
	iter, err := gitBranchIteratorNew(r.gitRepository, t)
	if err != nil {
		return nil, err
	}

	c := make(chan *Branch)
	w := &BranchWalker{
		gitBranchIterator: iter,
		C:                 c,
		co:                &sync.Once{},
		cc:                make(chan struct{}),
	}

	go w.run(c)
	return w, nil
}

// Cancel aborts an in-progress walk and drains the branch channel C.
func (w *BranchWalker) Cancel() {
	w.co.Do(w.cancel)
}

// Err returns error encountered while walking branches.
func (w *BranchWalker) Err() error {
	return w.err
}

// Slice returns a slice holding the branches and any error encountered while
// walking the branches.
func (w *BranchWalker) Slice() ([]*Branch, error) {
	s := []*Branch{}
	for b := range w.C {
		s = append(s, b)
	}
	return s, w.Err()
}

func (w *BranchWalker) cancel() {
	close(w.cc)
	for range w.C {
	}
}

func (w *BranchWalker) next() (*Branch, error) {
	r, t, err := w.gitBranchIterator.next()
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	return &Branch{r, t}, nil
}

func (w *BranchWalker) run(c chan<- *Branch) {
	defer close(c)

	for {
		branch, err := w.next()
		if err != nil {
			w.err = err
			return
		}
		if branch == nil {
			return
		}

		select {
		case c <- branch:
		case <-w.cc:
			return
		}
	}
}

type gitBranchIterator struct {
	ptr *C.git_branch_iterator
}

func (i *gitBranchIterator) init() {
	runtime.SetFinalizer(i, (*gitBranchIterator).free)
}

func (i *gitBranchIterator) free() {
	runtime.SetFinalizer(i, nil)
	C.git_branch_iterator_free(i.ptr)
}

func (i *gitBranchIterator) next() (*gitReference, branchType, error) {
	var (
		ptr *C.git_reference
		bt  C.git_branch_t
	)

	if err := unwrapErr(C.libgit2_branch_next(&ptr, &bt, i.ptr)); err != nil {
		return nil, 0, err
	}
	if ptr == nil {
		return nil, 0, nil
	}

	r := &gitReference{ptr}
	r.init()
	return r, branchType(bt), nil
}

func gitBranchCreate(repo *gitRepository, branchName string, target *gitCommit,
	force bool, signature *gitSignature, message string) (*gitReference, error) {

	var ptr *C.git_reference

	cname := C.CString(branchName)
	cforce := cbool(force)
	cmessage := C.CString(message)

	err := unwrapErr(C.libgit2_branch_create(&ptr, repo.ptr, cname, target.ptr,
		cforce, signature.ptr, cmessage))
	if err != nil {
		return nil, err
	}

	r := &gitReference{ptr}
	r.init()
	return r, nil
}

func gitBranchDelete(branch *gitReference) error {
	return unwrapErr(C.libgit2_branch_delete(branch.ptr))
}

func gitBranchIteratorNew(r *gitRepository, t branchType) (*gitBranchIterator, error) {
	var ptr *C.git_branch_iterator

	ct := C.git_branch_t(t)
	err := unwrapErr(C.libgit2_branch_iterator_new(&ptr, r.ptr, ct))
	if err != nil {
		return nil, err
	}

	i := &gitBranchIterator{ptr}
	i.init()
	return i, nil
}

func gitBranchName(r *gitReference) (string, error) {
	var cname *C.char
	defer C.free(unsafe.Pointer(cname))

	if err := unwrapErr(C.libgit2_branch_name(&cname, r.ptr)); err != nil {
		return "", err
	}
	return C.GoString(cname), nil
}
