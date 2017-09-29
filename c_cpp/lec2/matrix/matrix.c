#include <assert.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	assert(argc == 2);
	FILE *f = fopen(argv[1], "r");
	assert(f);

	int rows_n, cols_n;
	assert(fscanf(f, "%d %d",
				&rows_n, &cols_n) == 2);

	int **matrix = malloc(sizeof(int *) * rows_n);
	assert(matrix);

	for (int i = 0; i < rows_n; ++i) {
		matrix[i] = malloc(sizeof(int) * cols_n);
		assert(matrix[i]);
	}

	for (int i = 0; i < rows_n; ++i) {
		for (int j = 0; j < cols_n; ++j) {
			assert(fscanf(f, "%d",&matrix[i][j]) == 1);
		}
	}

	assert(!fclose(f));

	for (int i = 0; i < rows_n; ++i) {
		for (int j = 0; j < cols_n; ++j) {
			printf("%d ", matrix[i][j]);
		}
		printf("\n");
	}

	for (int i = 0; i < rows_n; ++i) {
		free(matrix[i]);
	}

	free(matrix);
}
