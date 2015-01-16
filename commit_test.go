package libgit2

import "testing"

func TestCreateEmptyCommit(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	_, err := repo.Commit(AllowEmptyMessage, AllowEmpty)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateCommit(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	tests := []struct {
		options []CommitOption
		err     error
	}{
		// clean working directory
		{err: errCommitIndexEmpty},

		// no error message
		{err: errCommitMessageEmpty, options: []CommitOption{AllowEmpty}},

		// valid empty commit
		{options: []CommitOption{AllowEmpty, AllowEmptyMessage}},
	}

	for _, test := range tests {
		_, err := repo.Commit(test.options...)
		if err != test.err {
			if err == nil {
				t.Errorf("want error %q, got none", err)
			} else {
				t.Errorf("want error %q, got %q", test.err, err)
			}
		}
	}
}
