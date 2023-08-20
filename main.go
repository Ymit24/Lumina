package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/kr/pretty"
)

func main() {
	lexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Ident", `[a-zA-Z]\w*`},
		{"Number", `[-+]?(\d*\.)?\d+`},
		{"Dot", `\.`},
		{"String", `\"(?:[^\"]|\\.)*\"`},
		{"Whitespace", `[ \t\n]+`},
		{"LParen", `\(`},
		{"RParen", `\)`},
		{"LBrace", `{`},
		{"RBrace", `}`},
		{"Colon", `:`},
		{"Comma", `,`},
		{"True", `true`},
		{"False", `false`},
		{"Semicolon", `;`},
	})
	parser := participle.MustBuild[Program](
		participle.Elide("Whitespace"),
		participle.Lexer(lexer),
	)

	bytes, err := os.ReadFile("first.la")
	if err != nil {
		panic(err)
	}

	contents := string(bytes)
	fmt.Printf("File: ```%s```\n", contents)

	program, err := parser.ParseString("", contents)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Query: %# v\n", pretty.Formatter(program))
	}
}
