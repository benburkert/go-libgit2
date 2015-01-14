package libgit2

//#include "libgit2.h"
import "C"

import (
	"os"
	"time"
	"unsafe"
)

type indexEntry struct {
	ptr *C.git_index_entry
}

// Base name of the file.
func (ie *indexEntry) Name() string {
	return C.GoString(ie.ptr.path)
}

// File length in bytes.
func (ie *indexEntry) Size() int64 {
	return int64(ie.ptr.file_size)
}

// File mode bits.
func (ie *indexEntry) Mode() os.FileMode {
	return os.FileMode(ie.ptr.mode)
}

// Modification time (mtime).
func (ie *indexEntry) ModTime() time.Time {
	return time.Unix(int64(ie.ptr.mtime.seconds), int64(ie.ptr.mtime.nanoseconds))
}

// IsDir returns true if the file is a directory.
func (ie *indexEntry) IsDir() bool {
	return ie.Mode().IsDir()
}

// Sys returns the underlying git_index_entry pointer.
func (ie *indexEntry) Sys() interface{} {
	return ie.ptr
}

var _ os.FileInfo = (*indexEntry)(nil)

func gitIndexGetBypath(idx *gitIndex, path string, stage int) *indexEntry {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ptr := C.git_index_get_bypath(idx.ptr, cpath, C.int(stage))
	if ptr == nil {
		return nil
	}
	return &indexEntry{ptr}
}
