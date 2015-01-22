package libgit2

type branchConfig struct {
	repo       Repository
	name       string
	target     *Commit
	force      bool
	sig        *Signature
	logMessage string
}

// BranchOption is an option type for git branch operations
type BranchOption func(*branchConfig)

func (c *branchConfig) check() error {
	var err error
	if c.target == nil {
		if c.target, err = c.repo.tip(); err != nil {
			return err
		}
	}

	if c.sig == nil {
		if c.sig, err = c.repo.DefaultSignature(); err != nil {
			return err
		}
	}

	return nil
}

// Target sets the commit to which a branch points.
func Target(target *Commit) BranchOption {
	return func(c *branchConfig) {
		c.target = target
	}
}

// Force overwrites an existing branch.
func Force() BranchOption {
	return func(c *branchConfig) {
		c.force = true
	}
}

// Creator sets the identity that will used to populate the reflog entry.
func Creator(sig *Signature) BranchOption {
	return func(c *branchConfig) {
		c.sig = sig
	}
}

// LogMessage is a one line long message to be appended to the reflog.
func LogMessage(message string) BranchOption {
	return func(c *branchConfig) {
		c.logMessage = message
	}
}
