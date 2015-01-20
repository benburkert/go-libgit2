package libgit2

import (
	"reflect"
	"testing"
)

func TestWalkRepo(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	n := 10
	mustSeedRepoN(t, repo, n)

	walk, err := repo.Walk()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for _ = range walk.C {
		count++
	}
	if err := walk.Err(); err != nil {
		t.Fatal(err)
	}

	if count != n {
		t.Errorf("want %d commits, got %d", n, count)
	}
}

func TestWalkerCancel(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	want, n := 7, 10
	mustSeedRepoN(t, repo, n)

	walk, err := repo.Walk()
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for _ = range walk.C {
		count++

		if count%want == 0 {
			walk.Cancel()
		}
	}

	if err := walk.Err(); err != nil {
		t.Fatal(err)
	}

	if count != want {
		t.Errorf("want %d commits, got %d", want, count)
	}
}

func TestWalkerBuffered(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	n := 10
	mustSeedRepoN(t, repo, n)

	walk, err := repo.Walk(BufferSize(3))
	if err != nil {
		t.Fatal(err)
	}

	for _ = range walk.C {
	}
	if err := walk.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestWalkerSlice(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	n := 10
	mustSeedRepoN(t, repo, n)

	walk, err := repo.Walk(BufferSize(3))
	if err != nil {
		t.Fatal(err)
	}

	commits, err := walk.Slice()
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != n {
		t.Errorf("want commits slice len %d, got %d", n, len(commits))
	}
}

func TestWalkerSortReverse(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	mustSeedRepoN(t, repo, 10)
	defer popd(t)

	walk, err := repo.Walk()
	if err != nil {
		t.Fatal(err)
	}
	wants, err := walk.Slice()
	if err != nil {
		t.Fatal(err)
	}

	rwalk, err := repo.Walk(Sorting(SortReverse))

	for i := len(wants) - 1; i > 0; i-- {
		want := wants[i]
		got := <-rwalk.C
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want commit %q during reverse walk, got %q", want, got)
		}
	}
}

func mustSeedRepo(t *testing.T, repo *Repository) {
	if _, err := repo.Commit(AllowEmpty, Message(rndstr()), CleanupMessage(false)); err != nil {
		t.Fatal(err)
	}
}

func mustSeedRepoN(t *testing.T, repo *Repository, n int) {
	for i := 0; i < n; i++ {
		mustSeedRepo(t, repo)
	}
}
