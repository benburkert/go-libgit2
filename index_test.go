package libgit2

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

func TestWriteIndex(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	idx, err := repo.Index()
	if err != nil {
		t.Fatal(err)
	}

	f := rndstr()
	if err = ioutil.WriteFile(f, []byte(rndstr()), 0644); err != nil {
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

func TestIndexWriteTree(t *testing.T) {
	repo := mustInitTestRepo(t)
	idx, err := repo.Index()
	if err != nil {
		t.Fatal(err)
	}

	f := mustSeedTestFile(t, repo)
	if err = idx.AddPath(f); err != nil {
		t.Fatal(err)
	}

	_, err = idx.WriteTree(*repo)
	if err != nil {
		t.Fatal(err)
	}
}

var (
	mu   sync.Mutex
	dirs = []string{}
)

func pushd(t *testing.T, dir string) {
	mu.Lock()
	defer mu.Unlock()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dirs = append(dirs, pwd)

	if err = os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}

func popd(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	dir := dirs[len(dirs)-1]
	dirs = dirs[:len(dirs)-1]

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}
