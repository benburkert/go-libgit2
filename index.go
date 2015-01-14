package libgit2

//#include "libgit2.h"
import "C"
import (
	"os"
	"runtime"
	"unsafe"
)

// Index is the in-memory representation of an index file.
type Index struct {
	*gitIndex
}

func repositoryIndex(repo Repository) (*Index, error) {
	i, err := gitRepositoryIndex(repo.gitRepository)
	if err != nil {
		return nil, err
	}
	return &Index{i}, nil
}

// AddPath adds a file by path to the index.
func (i Index) AddPath(path string) error {
	return gitIndexAddBypath(i.gitIndex, path)
}

// Get file info for a file in the index.
func (i Index) Get(path string) os.FileInfo {
	return gitIndexGetBypath(i.gitIndex, path, 0)
}

// Save the index on-disk.
func (i Index) Write() error {
	return gitIndexWrite(i.gitIndex)
}

type gitIndex struct {
	ptr *C.git_index
}

func (i *gitIndex) init() {
	runtime.SetFinalizer(i, (*gitIndex).free)
}

func (i *gitIndex) free() {
	runtime.SetFinalizer(i, nil)
	C.git_index_free(i.ptr)
}

func gitIndexAddBypath(idx *gitIndex, path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	return unwrapErr(C.libgit2_index_add_bypath(idx.ptr, cpath))
}

func gitIndexWrite(idx *gitIndex) error {
	return unwrapErr(C.libgit2_index_write(idx.ptr))
}

func gitRepositoryIndex(repo *gitRepository) (*gitIndex, error) {
	i := new(gitIndex)

	err := unwrapErr(C.libgit2_repository_index(&i.ptr, repo.ptr))
	if err != nil {
		return nil, err
	}
	i.init()
	return i, nil
}
