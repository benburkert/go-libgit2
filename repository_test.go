package libgit2

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRepositoryInitNormal(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

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

		want := pwd + "/" + test.path
		if got := repo.Path(); got != want {
			t.Errorf("got repo path %q, want %q", got, want)
		}

		want = pwd + "/" + test.workdir
		if got := repo.Workdir(); got != want {
			t.Errorf("got repo workdir %q, want %q", got, want)
		}

		if bare := repo.IsBare(); bare {
			t.Error("got bare repo, want normal")
		}
	}
}

func TestRepositoryInitBare(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct{ dir, path, workdir string }{
		{"testrepo.git/", "testrepo.git/", ""},
		{"noslash.git", "noslash.git/", ""},
	}

	for _, test := range tests {
		if _, err := os.Stat(test.dir); err == nil {
			t.Fatalf("%q directory exists", test.dir)
		} else if !os.IsNotExist(err) {
			t.Fatal(err)
		}

		repo, err := InitBareRepository(test.dir)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(test.path); err != nil {
			t.Fatal(err)
		}

		want := pwd + "/" + test.path
		if got := repo.Path(); got != want {
			t.Errorf("got repo path %q, want %q", got, want)
		}

		want = test.workdir
		if got := repo.Workdir(); got != want {
			t.Errorf("got repo workdir %q, want %q", got, want)
		}

		if bare := repo.IsBare(); !bare {
			t.Error("got normal repo, want bare")
		}
	}
}

func TestRepositoryDetachedHead(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	if repo.isDetachedHead() {
		t.Error("repo head is detached")
	}
}

func TestRepositoryUnbornHead(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	if !repo.isUnbornHead() {
		t.Error("repo head is not unborn")
	}
}

func TestRepositoryHead(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	_, err := repo.Head()
	gitErr, ok := err.(*gitError)
	if err == nil || !ok {
		t.Fatal("want errNotFound error")
	}
	if gitErr.class != errClassReference {
		t.Errorf("want error class %d, got %d", errClassNone, gitErr.class)
	}
	if gitErr.code != errUnbornBranch {
		t.Errorf("want error code %d, got %d", errUnbornBranch, gitErr.code)
	}
}

func mustInitTestRepo(t *testing.T) *Repository {
	repo, err := InitRepository(rndstr())
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

func mustSeedTestFile(t *testing.T, repo *Repository) string {
	pushd(t, repo.Workdir())
	defer popd(t)

	f := rndstr()
	if err := ioutil.WriteFile(f, []byte(rndstr()), 0644); err != nil {
		t.Fatal(err)
	}
	return f
}
