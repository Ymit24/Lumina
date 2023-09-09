build_and_run_example: build_lumina run_example
build_lumina:
	go build .
run_example:
	./lumina first.la first.ll
link_and_compile:
	llc -filetype=obj first.ll
	clang first.o -o executable -lm
run_bin:
	./executable
all: build_lumina run_example link_and_compile run_bin
