package libgit2

type walkerConfig struct {
	repo Repository

	startRef string
	bufSize  int
	sortMode SortMode
}

func (c *walkerConfig) check() error {
	return nil
}

// WalkerOption is an option type for walking operations.
type WalkerOption func(*walkerConfig)

// BufferSize sets the internal size of the commit channel for the walker.
func BufferSize(n int) WalkerOption {
	return func(c *walkerConfig) {
		c.bufSize = n
	}
}

// Sorting sets the sort mode of the walker.
func Sorting(mode SortMode) WalkerOption {
	return func(c *walkerConfig) {
		c.sortMode = mode
	}
}
