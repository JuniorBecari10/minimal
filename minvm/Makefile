FILES := main.c $(wildcard src/*.c)
FLAGS := -iquote include -std=c11 -Wall -Wextra
OUT := ../bin/minvm

all: debug release

debug:
	gcc $(FILES) $(FLAGS) -o $(OUT)d -g -pedantic -Wno-unused-function

release:
	gcc $(FILES) $(FLAGS) -o $(OUT)r -Werror -O3 -pedantic-errors

clean:
	-@rm -f $(OUT)-debug $(OUT)-release
