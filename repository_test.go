package libgit2

import (
	"os"
	"testing"
)

func TestRepositoryInit(t *testing.T) {
	tests := []struct{ dir, path, workdir string }{
		{"testrepo/", "testrepo/.git/", "testrepo/"},
		{"noslash", "noslash/.git/", "noslash/"},
	}

	for _, test := range tests {
		if _, err := os.Stat(test.dir); err == nil {
			t.Fatalf("%q directory exists", test.dir)
		} else if !os.IsNotExist(err) {
			t.Fatal(err)
		}

		repo, err := InitRepository(test.dir)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(test.path); err != nil {
			t.Fatal(err)
		}

		pwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		want := pwd + "/" + test.path
		if got := repo.Path(); got != want {
			t.Errorf("got repo path %q, want %q", got, want)
		}

		want = pwd + "/" + test.workdir
		if got := repo.Workdir(); got != want {
			t.Errorf("got repo workdir %q, want %q", got, want)
		}
	}
}
