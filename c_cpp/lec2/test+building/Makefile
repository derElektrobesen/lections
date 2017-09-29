all: binary

binary: main.o pass.o
	cc main.o pass.o -o binary

main.o: main.c
	cc -c main.c -o main.o

pass.o: pass.c
	cc -c pass.c -o pass.o
