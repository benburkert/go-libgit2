#include "libgit2.h"

libgit2_result libgit2_wrap_result(const int code)
{
       libgit2_result res = { code,NULL };

       if (res.code >= 0 || res.code == GIT_ITEROVER)
               return res;

       res.err = (git_error *)malloc(sizeof(git_error));
       giterr_detach(res.err);
       return res;
}

// branch.h

LIBGIT2_WRAPPER(libgit2_branch_create(
		git_reference **out,
		git_repository *repo,
		const char *branch_name,
		const git_commit *target,
		int force,
		const git_signature *signature,
		const char *log_message),
	git_branch_create(out, repo, branch_name, target, force, signature,
		log_message))

LIBGIT2_WRAPPER(libgit2_branch_iterator_new(
		git_branch_iterator **out,
		git_repository *repo,
		git_branch_t list_flags),
	git_branch_iterator_new(out, repo, list_flags))

LIBGIT2_WRAPPER(libgit2_branch_name(
		const char **out,
		const git_reference *ref),
	git_branch_name(out, ref))

LIBGIT2_WRAPPER(libgit2_branch_next(
		git_reference **out,
		git_branch_t *out_type,
		git_branch_iterator *iter),
	git_branch_next(out, out_type, iter))

// commit.h

LIBGIT2_WRAPPER(libgit2_commit_create(
		git_oid *id,
		git_repository *repo,
		const char *update_ref,
		const git_signature *author,
		const git_signature *committer,
		const char *message_encoding,
		const char *message,
		const git_tree *tree,
		size_t parent_count,
		const git_commit **parents),
	git_commit_create(id, repo, update_ref, author, committer, message_encoding,
		message, tree, parent_count, parents))

LIBGIT2_WRAPPER(libgit2_commit_lookup(
		git_commit **commit,
		git_repository *repo,
		const git_oid *id),
	git_commit_lookup(commit, repo, id))

LIBGIT2_WRAPPER(libgit2_commit_parent(
		git_commit **out,
		const git_commit *commit,
		unsigned int n),
	git_commit_parent(out, commit, n))

LIBGIT2_WRAPPER(libgit2_commit_parentcount(
		const git_commit *commit),
	git_commit_parentcount(commit))

// index.h

LIBGIT2_WRAPPER(libgit2_index_add_bypath(
		git_index *index,
		const char *path),
	git_index_add_bypath(index, path))

LIBGIT2_WRAPPER(libgit2_index_write(
		git_index *index),
	git_index_write(index))

LIBGIT2_WRAPPER(libgit2_index_write_tree(
		git_oid *out,
		git_index *index),
	git_index_write_tree(out, index))

// message.h

LIBGIT2_WRAPPER(libgit2_message_prettify(
		git_buf *out,
		const char *message,
		int strip_comments,
		char comment_char),
	git_message_prettify(out, message, strip_comments, comment_char))

// object.h

LIBGIT2_WRAPPER(libgit2_object_short_id(
		git_buf *out,
		const git_object *obj),
	git_object_short_id(out, obj))

// repository.h

LIBGIT2_WRAPPER(libgit2_repository_head(
		git_reference **out,
		git_repository *repo),
	git_repository_head(out, repo))

LIBGIT2_WRAPPER(libgit2_repository_head_detached(
		git_repository *repo),
	git_repository_head_detached(repo))

LIBGIT2_WRAPPER(libgit2_repository_index(
		git_index **out,
		git_repository *repo),
	git_repository_index(out, repo))

LIBGIT2_WRAPPER(libgit2_repository_init(
		git_repository **out,
		const char *path,
		unsigned int is_bare),
	git_repository_init(out, path, is_bare))

LIBGIT2_WRAPPER(libgit2_repository_open(
		git_repository **out,
		const char *path),
	git_repository_open(out, path))

// revwalk.h

LIBGIT2_WRAPPER(libgit2_revwalk_new(
		git_revwalk **out,
		git_repository *repo),
	git_revwalk_new(out, repo))

LIBGIT2_WRAPPER(libgit2_revwalk_next(
		git_oid *out,
		git_revwalk *walk),
	git_revwalk_next(out, walk))

LIBGIT2_WRAPPER(libgit2_revwalk_push_head(
		git_revwalk *walk),
	git_revwalk_push_head(walk))

// signature.h

LIBGIT2_WRAPPER(libgit2_signature_default(
		git_signature **out,
		git_repository *repo),
	git_signature_default(out, repo))

LIBGIT2_WRAPPER(libgit2_signature_dup(
		git_signature **dest,
		const git_signature *sig),
	git_signature_dup(dest, sig))

// tree.h

LIBGIT2_WRAPPER(libgit2_tree_lookup(
		git_tree **out,
		git_repository *repo,
		const git_oid *id),
	git_tree_lookup(out, repo, id))
