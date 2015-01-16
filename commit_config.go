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

	allowEmpty, allowEmptyMessage bool
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

	if c.message == "" && !c.allowEmptyMessage {
		return errCommitMessageEmpty
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

	return nil
}

// CommitOption is an option type for git commit options.
type CommitOption func(*commitConfig)

// AllowEmpty allows for an empty commit (one that has the exact same tree as
// its sole parent commit).
func AllowEmpty(c *commitConfig) { c.allowEmpty = true }

// AllowEmptyMessage allows for an empty commit message.
func AllowEmptyMessage(c *commitConfig) { c.allowEmptyMessage = true }
