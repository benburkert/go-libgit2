# go-libgit2 [![GoDoc](https://godoc.org/github.com/benburkert/go-libgit2?status.png)](http://godoc.org/github.com/benburkert/go-libgit2)

Idiomatic Go bindings for [libgit2](https://github.com/libgit2/libgit2).

## What's wrong with [git2go](https://github.com/libgit2/git2go)?

Short answer: nothing! git2go is great, you should probably be using it.

Long answer: git2go is a very thin layer over libgit2 so the Go API is nearly
identical to the C API. As a result, large parts of `git2go` don't feel very
Go-ish. This library is an experiment in implementing go binding for `libgit2`
in an idiomatic way. Compared to `git2go`, there are some big API changes, but
the internals look very similar.

## What makes it idiomatic?

### Walk commits via a channel or slice

The `Walker` type provides a `C` commit channel and `Slice()` function for
ranging over commits, instead of a C-like iterator type.

### Installable via `go get`

As long as you have a recent version of `libgit2` installed.

### No `runtime.LockOSThread`/`runtime.UnlockOSThread` calls required

Errors returned from `git_*` functions are wrapped in `libgit2_*` functions that
return the error & code in a result struct. Because this is done in C, there is
no need to manually lock the OS thread to retrieve the `libgit2` error.

### Value types

Most, if not all, of the exported types can be passed around as values. These
types hold internal pointers that manage the C pointers to `libgit2` objects,
so there's no need to worry about leaking memory or cleaning up after them.
