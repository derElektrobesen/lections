#include <stdbool.h>
#include <string.h>

bool is_pass_ok(const char *str) {
    return strcasecmp(str, "my_pass") == 0;
}
