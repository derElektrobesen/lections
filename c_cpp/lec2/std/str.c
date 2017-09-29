#include <stdio.h>
#include <assert.h>

int main(int argc, char *argv[]) {
    assert(argc == 2);
    char *buf = strdup(argv[1]);

    printf("buf='%s'\n",
            buf);
    free(buf);
    return 0;
}
