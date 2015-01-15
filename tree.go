package libgit2

//#include "libgit2.h"
import "C"

import "runtime"

type Tree struct {
	*gitTree
}

func lookupTree(repo Repository, oid OID) (*Tree, error) {
	tree, err := gitTreeLookup(repo.gitRepository, oid.gitOID)
	if err != nil {
		return nil, err
	}
	return &Tree{tree}, nil
}

type gitTree struct {
	ptr *C.git_tree
}

func (t *gitTree) init() {
	runtime.SetFinalizer(t, (*gitTree).free)
}

func (t *gitTree) free() {
	runtime.SetFinalizer(t, nil)
	C.git_tree_free(t.ptr)
}

func gitTreeLookup(repo *gitRepository, oid *gitOID) (*gitTree, error) {
	t := new(gitTree)

	err := unwrapErr(C.libgit2_tree_lookup(&t.ptr, repo.ptr, oid.ptr))
	if err != nil {
		return nil, err
	}
	t.init()
	return t, nil
}
