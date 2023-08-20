package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/kr/pretty"
)

type Query struct {
	Fields []*Field `@@*`
}

type Field struct {
	Source *Source `@@`
	Op     string  `@Comparison`
	Value  *Value  `@@`
}

type Source struct {
	Name string   `@Ident`
	Path []string `("." @Ident)*`
}

var queryLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Ident", `[a-zA-Z]\w*`},
	{"Number", `[-+]?(\d*\.)?\d+`},
	{"Comparison", `[=]|[<>]=?`},
	{"Dot", `\.`},
	{"String", `\"(?:[^\"]|\\.)*\"`},
	{"Whitespace", `[ \t]+`},
})

var _parser = participle.MustBuild[Query](
	participle.Unquote("String"),
	participle.Elide("Whitespace"),
	participle.Lexer(queryLexer),
)

func ParseQuery(q string) (*Query, error) {
	var expr *Query

	expr, err := _parser.ParseString("", q)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func Run() {
	q, err := ParseQuery(`a=100 item.price < 20 available >= 200`)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("Query: %# v", pretty.Formatter(q))
	}

	for _, field := range q.Fields {
		fmt.Printf("Field: '%s::%s': Value: %#v\n", field.Source.Name, field.Source.Path, field.Value)
	}
}
