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
