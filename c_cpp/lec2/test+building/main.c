#include <stdio.h>
#include "pass.h"

int main(int argc, char *argv[]) {
    (void)argc;
    printf("%s\n",
      is_pass_ok(argv[1]) ? "OK" : "FAIL");
}
