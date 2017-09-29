#include <stdio.h>
#include <assert.h>

int main(int argc, char *argv[]) {
    char buf[300];
    int n;
    assert(scanf("%d", &n) == 1);
    printf("scanned '%s', buf=%p\n",
      buf, buf);
    return 0;
}
