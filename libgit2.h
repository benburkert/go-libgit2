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

// branch.h

const libgit2_result libgit2_branch_create(
		git_reference **out,
		git_repository *repo,
		const char *branch_name,
		const git_commit *target,
		int force,
		const git_signature *signature,
		const char *log_message);

const libgit2_result libgit2_branch_delete(
		git_reference *branch);

const libgit2_result libgit2_branch_iterator_new(
		git_branch_iterator **out,
		git_repository *repo,
		git_branch_t list_flags);

const libgit2_result libgit2_branch_name(
		const char **out,
		const git_reference *ref);

const libgit2_result libgit2_branch_next(
		git_reference **out,
		git_branch_t *out_type,
		git_branch_iterator *iter);

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

const libgit2_result libgit2_commit_parent(
		git_commit **out,
		const git_commit *commit,
		unsigned int n);

const libgit2_result libgit2_commit_parentcount(
		const git_commit *commit);

// index.h

const libgit2_result libgit2_index_add_bypath(
		git_index *index,
		const char *path);

const libgit2_result libgit2_index_write(
		git_index *index);

const libgit2_result libgit2_index_write_tree(
		git_oid *out,
		git_index *index);

// message.h

const libgit2_result libgit2_message_prettify(
		git_buf *out,
		const char *message,
		int strip_comments,
		char comment_char);

// object.h

const libgit2_result libgit2_object_short_id(
		git_buf *out,
		const git_object *obj);

// repository.h

const libgit2_result libgit2_repository_head(
		git_reference **out,
		git_repository *repo);

const libgit2_result libgit2_repository_head_detached(
		git_repository *repo);

const libgit2_result libgit2_repository_index(
		git_index **out,
		git_repository *repo);

const libgit2_result libgit2_repository_init(
		git_repository **out,
		const char *path,
		unsigned int is_bare);

const libgit2_result libgit2_repository_open(
		git_repository **out,
		const char *path);

// revwalk.h

const libgit2_result libgit2_revwalk_new(
		git_revwalk **out,
		git_repository *repo);

const libgit2_result libgit2_revwalk_next(
		git_oid *out,
		git_revwalk *walk);

const libgit2_result libgit2_revwalk_push_head(
		git_revwalk *walk);

// signature.h

const libgit2_result libgit2_signature_default(
		git_signature **out,
		git_repository *repo);

const libgit2_result libgit2_signature_dup(
		git_signature **dest,
		const git_signature *sig);

// tree.h

const libgit2_result libgit2_tree_lookup(
		git_tree **out,
		git_repository *repo,
		const git_oid *id);

#endif
