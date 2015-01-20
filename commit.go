package libgit2

//#include "libgit2.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// Commit is the parsed representation of a commit object.
type Commit struct {
	*gitCommit
}

// Author returns the signature of the author of the commit.
func (c Commit) Author() (*Signature, error) {
	sig, err := c.author()
	if err != nil {
		return nil, err
	}
	return &Signature{sig}, nil
}

// Message is the full message of a commit.
func (c Commit) Message() string {
	return gitCommitMessage(c.gitCommit)
}

// ID is the object ID of the commit.
func (c Commit) ID() OID {
	return OID{gitCommitID(c.gitCommit)}
}

// Parents are the parent commits of the commit.
func (c Commit) Parents() ([]*Commit, error) {
	n, err := gitCommitParentcount(c.gitCommit)
	if err != nil {
		return nil, err
	}

	parents := make([]*Commit, n)
	for i := uint(0); i < n; i++ {
		cmt, err := gitCommitParent(c.gitCommit, i)
		if err != nil {
			return nil, err
		}
		parents[i] = &Commit{cmt}
	}
	return parents, nil
}

// ShortID returns an abbreviated object ID of the commit.
func (c Commit) ShortID() (string, error) {
	return c.shortID()
}

func (c Commit) String() string {
	return c.ID().String()
}

func createCommit(config *commitConfig) (*Commit, error) {
	gitParents := make([]*gitCommit, len(config.parents))
	for i, c := range config.parents {
		gitParents[i] = c.gitCommit
	}

	oid, err := gitCommitCreate(config.repo.gitRepository, config.updateRef,
		config.author.gitSignature, config.committer.gitSignature, config.encoding,
		config.message, config.tree.gitTree, gitParents)
	if err != nil {
		return nil, err
	}
	return lookupCommit(config.repo, OID{oid})
}

func lookupCommit(repo Repository, oid OID) (*Commit, error) {
	cmt, err := gitCommitLookup(repo.gitRepository, oid.gitOID)
	if err != nil {
		return nil, err
	}
	return &Commit{cmt}, nil
}

type gitCommit struct {
	ptr *C.git_commit
}

func (c *gitCommit) author() (*gitSignature, error) {
	return gitCommitAuthor(c).dup()
}

func (c *gitCommit) init() {
	runtime.SetFinalizer(c, (*gitCommit).free)
}

func (c *gitCommit) free() {
	runtime.SetFinalizer(c, nil)
	C.git_commit_free(c.ptr)
}

func (c *gitCommit) shortID() (string, error) {
	buf := &C.git_buf{}
	defer C.git_buf_free(buf)

	err := unwrapErr(C.libgit2_object_short_id(buf, (*C.git_object)(c.ptr)))
	if err != nil {
		return "", err
	}
	return C.GoString(buf.ptr), nil
}

func gitCommitAuthor(commit *gitCommit) *gitSignature {
	return &gitSignature{ptr: C.git_commit_author(commit.ptr)}
}

func gitCommitCreate(repo *gitRepository, updateRef string, author,
	committer *gitSignature, messageEncoding, message string, tree *gitTree,
	parents []*gitCommit) (*gitOID, error) {

	oid := &gitOID{ptr: &C.git_oid{}}

	var cref *C.char
	if updateRef != "" {
		cref = C.CString(updateRef)
		defer C.free(unsafe.Pointer(cref))
	}

	var cenc *C.char
	if messageEncoding != "" {
		cenc = C.CString(messageEncoding)
		defer C.free(unsafe.Pointer(cenc))
	}

	cmsg := C.CString(message)
	defer C.free(unsafe.Pointer(cmsg))

	var cparents **C.git_commit
	if len(parents) > 0 {
		ary := make([]*C.git_commit, len(parents))
		for i, v := range parents {
			ary[i] = v.ptr
		}
		cparents = &ary[0]
	}

	return oid, unwrapErr(C.libgit2_commit_create(oid.ptr, repo.ptr, cref,
		author.ptr, committer.ptr, cenc, cmsg, tree.ptr, C.size_t(len(parents)),
		cparents))
}

func gitCommitLookup(repo *gitRepository, oid *gitOID) (*gitCommit, error) {
	c := new(gitCommit)

	err := unwrapErr(C.libgit2_commit_lookup(&c.ptr, repo.ptr, oid.ptr))
	if err != nil {
		return nil, err
	}
	c.init()
	return c, nil
}

func gitCommitID(commit *gitCommit) *gitOID {
	return &gitOID{C.git_commit_id(commit.ptr)}
}

func gitCommitMessage(commit *gitCommit) string {
	return C.GoString(C.git_commit_message(commit.ptr))
}

func gitCommitParent(commit *gitCommit, n uint) (*gitCommit, error) {
	c := new(gitCommit)
	return c, unwrapErr(C.libgit2_commit_parent(&c.ptr, commit.ptr, C.uint(n)))
}

func gitCommitParentcount(commit *gitCommit) (uint, error) {
	res := C.libgit2_commit_parentcount(commit.ptr)
	return uint(res.code), unwrapErr(res)
}
