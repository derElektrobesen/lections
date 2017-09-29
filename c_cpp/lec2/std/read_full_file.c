#include <assert.h>
#include <sys/stat.h>
#include <stdio.h>
#include <stdlib.h>

static size_t get_file_size(
  const char *path) {
    struct stat st;
    if (lstat(path, &st))
        abort();

    return st.st_size;
}

static const char *read_full_file(
  const char *path) {
    size_t file_size = get_file_size(path);
    printf("file size = %zu\n", file_size);
    assert(file_size);
    char *contents = malloc(file_size);
    assert(contents);

    FILE *f = fopen(path, "r");
    int num;
    int read_bytes = fread(contents, 1,
      file_size, f);
    assert(read_bytes == file_size);
    fclose(f);

    return contents;
}

int main(int argc, char *argv[]) {
    const char *contents =
      read_full_file(argv[1]);
    printf("file contents = '%s'\n",
      contents);
    free((void *)contents);

    return 0;
}

// sudo dtruss ./a.out test.txt 2>&1 | fgrep -A100 "test.txt"
