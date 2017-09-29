#include <assert.h>
#include <stdlib.h>

#include "pass.h"

int main() {
    assert(is_pass_ok("") == false);
    assert(is_pass_ok(NULL) == false);
    assert(is_pass_ok("my_pass") == true);
    assert(is_pass_ok("My_pass") == true);
    assert(is_pass_ok("invalid") == false);
}
