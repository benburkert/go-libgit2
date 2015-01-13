package libgit2

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	dir := setup()
	r := m.Run()
	cleanup(dir)
	os.Exit(r)
}

func setup() string {
	dir, err := ioutil.TempDir("", "go-libgit2")
	if err != nil {
		panic(err)
	}

	if err = os.Chdir(dir); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(".gitconfig", gitconfig, 0644); err != nil {
		panic(err)
	}
	os.Setenv("HOME", dir)

	Init()
	return dir
}

func cleanup(dir string) {
	return
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
	Shutdown()
}

var gitconfig = []byte(`
[user]
  name = Default
  email = default@example.com
`)
