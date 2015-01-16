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
		message string
	}{
		// clean working directory
		{err: errCommitIndexEmpty},

		// no error message
		{err: errCommitMessageEmpty, options: []CommitOption{AllowEmpty}},

		// empty commit, empty message
		{options: []CommitOption{AllowEmpty, AllowEmptyMessage}},

		// untouched message
		{
			options: []CommitOption{AllowEmpty, Message("untouched message")},
			message: "untouched message",
		},

		// newline added to message
		{
			options: []CommitOption{
				AllowEmpty,
				Message("newline added"),
				CleanupMessage(false),
			},
			message: "newline added\n",
		},

		// newline added, comment untouched
		{
			options: []CommitOption{
				AllowEmpty,
				Message("newline added\n#comment untouched"),
				CleanupMessage(false),
			},
			message: "newline added\n#comment untouched\n",
		},

		// comment stripped
		{
			options: []CommitOption{
				AllowEmpty,
				Message("comment stripped\n#stripped comment\n"),
				CleanupMessage(true),
			},
			message: "comment stripped\n",
		},

		// comment with custom marker stripped
		{
			options: []CommitOption{
				AllowEmpty,
				Message("weird comment\n!this is a comment line\n"),
				CleanupMessage(true),
				CommentMarker('!'),
			},
			message: "weird comment\n",
		},
	}

	for _, test := range tests {
		comment, err := repo.Commit(test.options...)
		if err != test.err {
			if err == nil {
				t.Errorf("want error %q, got none", err)
			} else {
				t.Errorf("want error %q, got %q", test.err, err)
			}
		}

		if test.message != "" && test.message != comment.Message() {
			t.Errorf("want comment message %q, got %q", test.message, comment.Message())
		}
	}
}
