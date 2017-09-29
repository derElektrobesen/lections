#include <stdio.h>

static void read_file(const char *path) {
    FILE *f = fopen(path, "r");
    int num;
    while (fscanf(f, "%d", &num) == 1) {
        printf("read num = %d\n", num);
    }
    fclose(f);
}

int main(int argc, char *argv[]) {
    read_file(argv[1]);
    return 0;
}

//sudo dtruss ./a.out test.txt 2>&1 | fgrep -A100 "test.txt"
