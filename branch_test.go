package libgit2

import "testing"

func TestCreateBranch(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	mustSeedRepo(t, repo)

	branches := map[string]*Branch{}
	for i := 0; i < 10; i++ {
		branch, err := repo.CreateBranch(rndstr())
		if err != nil {
			t.Fatal(err)
		}
		name, err := branch.Name()
		if err != nil {
			t.Fatal(err)
		}
		branches[name] = branch
	}

	walker, err := repo.Branches()
	if err != nil {
		t.Fatal(err)
	}

	for branch := range walker.C {
		want, err := branch.Name()
		if err != nil {
			t.Fatal(err)
		}
		if want == "master" {
			continue
		}

		branch, ok := branches[want]
		if !ok {
			t.Fatalf("unexpected branch %q", want)
		}

		got, err := branch.Name()
		if err != nil {
			t.Fatal(err)
		}

		if want != got {
			t.Errorf("want branch name %q, got %q", want, got)
		}
	}
}

func TestDestroyBranch(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	mustSeedRepo(t, repo)

	branch, err := repo.CreateBranch(rndstr())
	if err != nil {
		t.Fatal(err)
	}

	if err = branch.Delete(); err != nil {
		t.Error(err)
	}
}

func TestMoveBranch(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	mustSeedRepo(t, repo)

	oldBranch, err := repo.CreateBranch(rndstr())
	if err != nil {
		t.Fatal(err)
	}

	want := rndstr()
	newBranch, err := oldBranch.Move(want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := newBranch.Name()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want branch name %q, got %q", want, got)
	}
}

func TestRenameBranch(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	mustSeedRepo(t, repo)

	branch, err := repo.CreateBranch(rndstr())
	if err != nil {
		t.Fatal(err)
	}

	want := rndstr()
	if err := branch.Rename(want); err != nil {
		t.Fatal(err)
	}

	got, err := branch.Name()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want branch name %q, got %q", want, got)
	}
}
