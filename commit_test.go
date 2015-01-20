package libgit2

import (
	"reflect"
	"strings"
	"testing"
)

func TestCreateEmptyCommit(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	if _, err := repo.Commit(AllowEmptyMessage, AllowEmpty); err != nil {
		t.Fatal(err)
	}
}

func TestCreateCommitMessage(t *testing.T) {
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
			} else if test.err == nil {
				t.Error(err)
			} else {
				t.Errorf("want error %q, got %q", test.err, err)
			}
			break
		}

		if test.message != "" && test.message != comment.Message() {
			t.Errorf("want comment message %q, got %q", test.message, comment.Message())
		}
	}
}

func TestCommitAuthor(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	want, err := repo.DefaultSignature()
	if err != nil {
		t.Fatal(err)
	}

	commit, err := repo.Commit(Message("testing commit author"), AllowEmpty)
	if err != nil {
		t.Fatal(err)
	}

	got, err := commit.Author()
	if err != nil {
		t.Fatal(err)
	}

	if want.Name != got.Name {
		t.Errorf("want author name %q, got %q", want.Name, got.Name)
	}

	if want.Email != got.Email {
		t.Errorf("want author email %q, got %q", want.Email, got.Email)
	}

	if !want.When.Equal(got.When) {
		t.Errorf("want author timestamp %q, got %q", want.When, got.When)
	}
}

func TestCommitParents(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	parent, err := repo.Commit(Message("parent"), AllowEmpty)
	if err != nil {
		t.Fatal(err)
	}

	commit, err := repo.Commit(
		Message("child"),
		Parents(parent),
		AllowEmpty)
	if err != nil {
		t.Fatal(err)
	}

	want := []*Commit{parent}
	got, err := commit.Parents()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want commit parents %v, got %v", want, got)
	}
}

func TestCommitShortID(t *testing.T) {
	repo := mustInitTestRepo(t)
	pushd(t, repo.Workdir())
	defer popd(t)

	commit, err := repo.Commit(AllowEmpty, AllowEmptyMessage)
	if err != nil {
		t.Fatal(err)
	}

	shortID, err := commit.ShortID()
	if err != nil {
		t.Fatal(err)
	}

	id := commit.ID().String()
	if !strings.HasPrefix(id, shortID) {
		t.Errorf("invalid short id %q for %q", shortID, id)
	}
}
