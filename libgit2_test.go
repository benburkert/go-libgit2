package libgit2

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "go-libgit2")
	if err != nil {
		panic(err)
	}

	if err = os.Chdir(dir); err != nil {
		panic(err)
	}

	Init()
	r := m.Run()
	if err = os.RemoveAll(dir); err != nil {
		panic(err)
	}
	Shutdown()
	os.Exit(r)
}
