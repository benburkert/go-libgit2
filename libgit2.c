#include "libgit2.h"

libgit2_result libgit2_wrap_result(const int code)
{
       libgit2_result res = { code,NULL };

       if (res.code >= 0)
               return res;

       res.err = (git_error *)malloc(sizeof(git_error));
       giterr_detach(res.err);
       return res;
}

LIBGIT2_WRAPPER(libgit2_repository_init(
               git_repository **out,
               const char *path,
               unsigned int is_bare),
       git_repository_init(out, path, is_bare))
