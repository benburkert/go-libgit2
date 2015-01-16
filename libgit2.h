#ifndef _LIBGIT2_H_
#define _LIBGIT2_H_

#include <git2.h>

#define LIBGIT2_WRAPPER(sig, call) \
const libgit2_result sig { \
       return libgit2_wrap_result(call); \
}

typedef struct libgit2_result {
       int code;
       git_error *err;
} libgit2_result;

libgit2_result libgit2_wrap_result(const int);

// commit.h

const libgit2_result libgit2_commit_create(
		git_oid *id,
		git_repository *repo,
		const char *update_ref,
		const git_signature *author,
		const git_signature *committer,
		const char *message_encoding,
		const char *message,
		const git_tree *tree,
		size_t parent_count,
		const git_commit **parents);

const libgit2_result libgit2_commit_lookup(
		git_commit **commit,
		git_repository *repo,
		const git_oid *id);

// index.h

const libgit2_result libgit2_index_add_bypath(
		git_index *index,
		const char *path);

const libgit2_result libgit2_index_write(
		git_index *index);

const libgit2_result libgit2_index_write_tree(
		git_oid *out, git_index *index);

// repository.h

const libgit2_result libgit2_repository_index(
		git_index **out,
		git_repository *repo);

const libgit2_result libgit2_repository_init(
		git_repository **out,
		const char *path,
		unsigned int is_bare);

// signature.h

const libgit2_result libgit2_signature_default(
		git_signature **out,
		git_repository *repo);

// tree.h

const libgit2_result libgit2_tree_lookup(
		git_tree **out,
		git_repository *repo,
		const git_oid *id);

#endif
