package libgit2

import (
	"io/ioutil"
	"testing"
)

func TestWriteIndex(t *testing.T) {
	repo := mustInitTestRepo(t)
	idx, err := repo.Index()
	if err != nil {
		t.Fatal(err)
	}

	f := rndstr()
	if ioutil.WriteFile(f, []byte(rndstr()), 0644); err != nil {
		t.Fatal(err)
	}

	if err = idx.AddPath(f); err != nil {
		t.Fatal(err)
	}

	if err = idx.Write(); err != nil {
		t.Fatal(err)
	}

	// refetch index
	if idx, err = repo.Index(); err != nil {
		t.Fatal(err)
	}

	info := idx.Get(f)
	if f != info.Name() {
		t.Errorf("want index file name %q, got %q", f, info.Name())
	}
}
