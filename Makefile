build_and_run_example: build_lumina run_example
build_lumina:
	go build .
	mv lumina out/
run_example:
	./out/lumina first.la out/first.ll
link_and_compile:
	llc -filetype=obj out/first.ll
	clang out/first.o -o out/executable -lm
run_bin:
	./out/executable
all: build_lumina run_example link_and_compile run_bin
