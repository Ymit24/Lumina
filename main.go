package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/kr/pretty"
)

type Query struct {
	Fields []*Field `@@*`
}

type Field struct {
	Key   string `@Ident "="`
	Op    string `@("=" | "<" "=" | "<" | ">" "=" | ">" )`
	Value *Value `@@`
}

type Value struct {
	Number *float64 ` @Float | @Int`
	String *string  `| @String`
	Bool   *bool    `| ( @"true" | @"false" )`
	Nil    bool     `| @"nil"`
}

var parser = participle.MustBuild[Query](
	participle.Unquote("String"),
)

func ParseQuery(q string) (*Query, error) {
	var expr *Query

	expr, err := parser.ParseString("", q)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func main() {
	q, err := ParseQuery(`a=100 price < 20 available >= 200`)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("Query: %# v", pretty.Formatter(q))
	}

	for _, field := range q.Fields {
		fmt.Printf("Field: '%s': Value: %#v\n", field.Key, field.Value)
	}
}
