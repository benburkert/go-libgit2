package libgit2

//#include "libgit2.h"
import "C"

import (
	"runtime"
	"unsafe"
)

// Repository is an on-disk Git repository.
type Repository struct {
	*gitRepository
}

// InitBareRepository initializes a are Git repository.
func InitBareRepository(dir string) (*Repository, error) {
	r, err := gitInitRepository(dir, true)
	if err != nil {
		return nil, err
	}

	return &Repository{r}, nil
}

// InitRepository initializes a normal Git repository.
func InitRepository(dir string) (*Repository, error) {
	r, err := gitInitRepository(dir, false)
	if err != nil {
		return nil, err
	}

	return &Repository{r}, nil
}

func (r Repository) Commit(options ...CommitOption) (*Commit, error) {
	config := &commitConfig{repo: r}
	for _, opt := range options {
		opt(config)
	}
	if err := config.check(); err != nil {
		return nil, err
	}

	return createCommit(config)
}

// DefaultSignature returns a new action signature with default user and now
// timestamp.
func (r Repository) DefaultSignature() (*Signature, error) {
	return defaultSignature(r)
}

// Retrieve and resolve the reference pointed at by HEAD.
func (r Repository) Head() (*Reference, error) {
	ref, err := gitRepositoryHead(r.gitRepository)
	if err != nil {
		return nil, err
	}
	return &Reference{ref}, nil
}

// Index returns the index file for the repository.
func (r Repository) Index() (*Index, error) {
	return repositoryIndex(r)
}

// IsBare returns true if the repository is does not contain a working
// directory.
func (r Repository) IsBare() bool {
	return gitRepositoryIsBare(r.gitRepository)
}

// Path returns the file path the .git directory for normal repositories, or
// the repository itself for bare repositories.
func (r Repository) Path() string {
	return gitRepositoryPath(r.gitRepository)
}

// Workdir returns the file path of the working directory for the repository.
func (r Repository) Workdir() string {
	return gitRepositoryWorkdir(r.gitRepository)
}

func (r Repository) isDetachedHead() bool {
	return gitRepositoryHeadDetached(r.gitRepository)
}

func (r Repository) isUnbornHead() bool {
	return gitRepositoryHeadUnborn(r.gitRepository)
}

func (r Repository) tip() (*Commit, error) {
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}
	return lookupCommit(r, *ref.target())
}

type gitRepository struct {
	ptr *C.git_repository
}

func (r *gitRepository) init() {
	runtime.SetFinalizer(r, (*gitRepository).free)
}

func (r *gitRepository) free() {
	runtime.SetFinalizer(r, nil)
	C.git_repository_free(r.ptr)
}

func gitInitRepository(path string, isBare bool) (*gitRepository, error) {
	r := new(gitRepository)

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cbare := ucbool(isBare)

	if err := unwrapErr(C.libgit2_repository_init(&r.ptr, cpath, cbare)); err != nil {
		return nil, err
	}
	r.init()
	return r, nil
}

func gitRepositoryHead(repo *gitRepository) (*gitReference, error) {
	r := new(gitReference)

	if err := unwrapErr(C.libgit2_repository_head(&r.ptr, repo.ptr)); err != nil {
		return nil, err
	}
	r.init()
	return r, nil
}

func gitRepositoryHeadDetached(repo *gitRepository) bool {
	return C.git_repository_head_detached(repo.ptr) != 0
}

func gitRepositoryHeadUnborn(repo *gitRepository) bool {
	return C.git_repository_head_unborn(repo.ptr) != 0
}

func gitRepositoryIsBare(repo *gitRepository) bool {
	return C.git_repository_is_bare(repo.ptr) != 0
}

func gitRepositoryPath(repo *gitRepository) string {
	return C.GoString(C.git_repository_path(repo.ptr))
}

func gitRepositoryWorkdir(repo *gitRepository) string {
	return C.GoString(C.git_repository_workdir(repo.ptr))
}
