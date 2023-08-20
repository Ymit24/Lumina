build_and_run_example: build run_example
build:
	go build -o lumina main.go grammar.go
run_example:
	./lumina first.la
