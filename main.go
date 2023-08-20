package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/kr/pretty"
)

func main() {
	parser := participle.MustBuild[Program](
		participle.Elide("Whitespace"),
		participle.Lexer(luminaLexer),
	)

	if len(os.Args) != 2 {
		panic("Expected lumina file to be passed in!")
	}
	filename := os.Args[1]

	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	contents := string(bytes)
	fmt.Printf("File: ```%s```\n", contents)

	program, err := parser.ParseString(filename, contents, participle.Trace(os.Stdout))

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Query: %# v\n", pretty.Formatter(program))
	}
}
