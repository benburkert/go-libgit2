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

// repository.h

const libgit2_result libgit2_repository_init(
		git_repository **out,
		const char *path,
		unsigned int is_bare);

// signature.h

const libgit2_result libgit2_signature_default(
		git_signature **out,
		git_repository *repo);

#endif
