FILES=main.c src/*.c
FLAGS=-Iinclude -std=c11 -Wall -Wextra

all: debug

debug:
	gcc $(FILES) $(FLAGS) -g -pedantic

release:
	gcc $(FILES) $(FLAGS) -Werror -O3 -pedantic-errors

