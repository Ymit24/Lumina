package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/kr/pretty"
)

func main() {
	fmt.Println("Lumina Compiler")
	parser := participle.MustBuild[Program](
		participle.Elide("Whitespace"),
		participle.Lexer(luminaLexer),
	)

	if len(os.Args) != 3 {
		panic("Expected lumina file to be passed in and output binary name!")
	}
	inputFilename := os.Args[1]
	outputFilename := os.Args[2]

	fmt.Println("Reading source file")
	bytes, err := os.ReadFile(inputFilename)
	if err != nil {
		panic(err)
	}

	contents := string(bytes)
	fmt.Printf("File: ```%s```\n", contents)

	program, err := parser.ParseString(inputFilename, contents, participle.Trace(os.Stdout))

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Query: %# v\n", pretty.Formatter(program))
	}

	fmt.Println("Visiting")

	codeGenerator := NewCodeGenerator()

	code := codeGenerator.VisitProgram(program)

	os.WriteFile(
		outputFilename,
		[]byte(code),
		0644,
	)
}
