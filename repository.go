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

// InitRepository initializes a normal Git repository.
func InitRepository(dir string) (*Repository, error) {
	r, err := gitInitRepository(dir, false)
	if err != nil {
		return nil, err
	}

	return &Repository{r}, nil
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

type gitRepository struct {
	ptr *C.git_repository
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
	runtime.SetFinalizer(r, (*gitRepository).free)
	return r, nil
}

func gitRepositoryPath(repo *gitRepository) string {
	return C.GoString(C.git_repository_path(repo.ptr))
}

func gitRepositoryWorkdir(repo *gitRepository) string {
	return C.GoString(C.git_repository_workdir(repo.ptr))
}
