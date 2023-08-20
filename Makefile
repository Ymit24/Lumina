build_and_run_example: build run_example
build:
	go build -o lumina main.go grammar.go codegen.go types.go
run_example:
	./lumina first.la
