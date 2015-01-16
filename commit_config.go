package libgit2

import "errors"

var (
	errCommitIndexEmpty   = errors.New("nothing to commit, working directory clean")
	errCommitMessageEmpty = errors.New("empty commit message")
)

type commitConfig struct {
	repo Repository

	author, committer *Signature

	index   *Index
	tree    *Tree
	parents []*Commit

	encoding, message, updateRef string

	allowEmpty, allowEmptyMessage, allowOrphan bool

	cleanupMessage, stripComments bool
	commentMarker                 rune
}

func (c *commitConfig) check() error {
	if c.tree == nil {
		if c.index == nil {
			idx, err := c.repo.Index()
			if err != nil {
				return err
			}
			if idx.entryCount() == 0 && !c.allowEmpty {
				return errCommitIndexEmpty
			}
			c.index = idx
		}

		tree, err := c.index.WriteTree(c.repo)
		if err != nil {
			return err
		}
		c.tree = tree
	}

	if c.updateRef == "" {
		c.updateRef = "HEAD"
	}

	if c.message == "" && !c.allowEmptyMessage {
		return errCommitMessageEmpty
	}

	if c.cleanupMessage {
		msg, err := gitMessagePrettify(c.message, c.stripComments, c.commentMarker)
		if err != nil {
			return err
		}
		c.message = msg
	}

	if c.author == nil {
		sig, err := c.repo.DefaultSignature()
		if err != nil {
			return err
		}
		c.author = sig
	}

	if c.committer == nil {
		c.committer = c.author
	}

	if c.repo.isUnbornHead() {
		c.allowOrphan = true
	}

	if len(c.parents) == 0 && !c.allowOrphan {
		tip, err := c.repo.tip()
		if err != nil {
			return err
		}
		c.parents = []*Commit{tip}
	}

	return nil
}

// CommitOption is an option type for git commit options.
type CommitOption func(*commitConfig)

// AllowEmpty allows for an empty commit (one that has the exact same tree as
// its sole parent commit).
func AllowEmpty(c *commitConfig) { c.allowEmpty = true }

// AllowEmptyMessage allows for an empty commit message.
func AllowEmptyMessage(c *commitConfig) { c.allowEmptyMessage = true }

// AllowOrphan allows for an orphaned commit to be created.
func AllowOrphan(c *commitConfig) { c.allowOrphan = true }

// CleanupMessage automatically strips whitespace and adds a newline at the end
// of the commit message. If stripComments is true, comment lines are removed.
func CleanupMessage(stripComments bool) CommitOption {
	return func(c *commitConfig) {
		c.cleanupMessage = true
		if stripComments {
			c.stripComments = true
			if c.commentMarker == 0 {
				c.commentMarker = '#'
			}
		}
	}
}

// CommentMarker sets the rune character of comment lines in the message.
func CommentMarker(char rune) CommitOption {
	return func(c *commitConfig) {
		c.commentMarker = char
	}
}

// Message sets the commit message string.
func Message(message string) CommitOption {
	return func(c *commitConfig) {
		c.message = message
	}
}
