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

#endif
