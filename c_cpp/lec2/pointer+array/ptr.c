#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

static void f1() {
	int a = 1;
	printf("&a = %lld\n", (unsigned long long)&a);
	printf("&a = %p\n", &a);
}

static void f2() {
	char a = 1;
	char *b = &a;
	printf("a = %p\n", &a);
	printf("b = %p\n", b);
}

static void f3() {
	int a = 1;
	int *b = &a;
	printf("*b = %d\n", *b);

	*b += 1;
	printf("*b = %d\n", *b);
	printf("a = %d\n", a);
}

static void f4() {
	int a = 1;
	int *b = &a;
	int **c = &b;
	**c = 3;
	printf("a = %d\n", a);
}

static void swap(int x, int y) {
	int temp;
	temp = x;
	x = y;
	y = temp;
}
static void swap2(int *px, int *py) {
	int temp;
	temp = *px;
	*px = *py;
	*py = temp;
}
static void f5() {
	int a = 1, b = 2;
	swap2(&a, &b);
	printf("a=%d, b=%d\n", a, b);
}

static void f6() {
	int a[] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
	printf("a[2] = %d\n", a[2]);
	printf("a[10] = %d\n", a[10]);
}

static void f7() {
	int a[10] = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
	int *p1 = &a[3];
	int *p2 = a + 3;
	assert(p1 == p2);

	assert(a[3] == *(a + 3));

	assert(a[4] == 4[a]);
	//assert(*(a+4) == *(4+a));

	assert(*(a + 4) == *(4 + a));
}

static size_t calc_string_len(
		const char *s) {

	const char *p = s;
	while (*p != '\0')
		++p;
	return p - s;
}

static void f8() {
	const char a[] = "hello";
	printf("a[4] = %c\n", a[4]);
	printf("a[5] = '%c' (%d)\n", a[5], a[5]);
	char *p = (char *)&a[4];
	//a[4] = 'a';
	//printf("a=%s\n", a_copy);
	//    printf("len(a)=%zu\n", strlen(a));
}

static void f8_1() {
	const char *s = "abc";
	printf("mystrlen(s) = %zu\n", strlen(s));
}

static void my_strcpy(char *s, char *t) {
	while ((*s++ = *t++) != '\0');
	//  while (*t != '\0') {
	//    *s++ = *t++;
	//}
	//*s = '\0';
}

int main(int argc, char *argv[]) {
	f1();
}
