package main

import "github.com/alecthomas/participle/v2/lexer"

type Program struct {
	Pos        lexer.Position
	Statements []*HighLevelStatement `@@*`
}

type HighLevelStatement struct {
	Pos      lexer.Position
	Extern   *Extern   `@@`
	Function *Function `| @@`
}

type Extern struct {
	Pos       lexer.Position
	Signiture FunctionSignature `Extern @@ Semicolon`
}

type FunctionSignature struct {
	Pos        lexer.Position
	Name       string              `Fn @Ident`
	Args       []*FunctionArgument `( LParen ( @@ ( Comma @@ )* )? RParen )`
	ReturnType *Type               `(Colon @@)?`
}

type Function struct {
	Pos       lexer.Position
	Signature FunctionSignature `@@`
	Block     CodeBlock         `@@`
}
type Type struct {
	Pos   lexer.Position
	Array *ArrayType `@@`
	Inner *InnerType `| @@`
}
type InnerType struct {
	Pos              lexer.Position
	Name             string  `@Ident`
	GenericArguments []*Type `( LParen ( @@ ( Comma @@ )* ) RParen )?`
}
type ArrayType struct {
	Pos      lexer.Position
	IsSpread bool `@Ellipsis?`
	Type     Type `LBracket @@ RBracket`
}

type CodeBlock struct {
	Pos        lexer.Position
	Statements []Statement `LBrace ( @@ )* RBrace`
}
type FunctionArgument struct {
	Pos  lexer.Position
	Name string `@Ident`
	Type Type   `Colon @@`
}
type Statement struct {
	Pos                lexer.Position
	VariableAssignment *VariableAssignment `( @@`
	Return             *Return             `| @@`
	FunctionCall       *FunctionCall       `| @@ ) Semicolon`
}
type Return struct {
	Expression *Expression `Return (@@)?`
}
type VariableAssignment struct {
	Pos        lexer.Position
	Mutability string     `@(Static | Const | Var)`
	Name       string     `@Ident`
	Type       *Type      `Colon ( @@ )?`
	Expression Expression `Equals @@`
}
type FunctionCall struct {
	Pos          lexer.Position
	FunctionName string       `@Ident`
	Args         []Expression `LParen ( @@ (Comma @@)* )? RParen`
}
type Expression struct {
	Pos    lexer.Position
	Term   Term        `@@`
	AddSub *string     `[ @(Plus | Minus)`
	Next   *Expression `@@ ]`
}
type Term struct {
	Pos    lexer.Position
	Factor Value   `@@`
	MulDiv *string `[ @(Asterisk | FSlash)`
	Next   *Term   `@@ ]`
}
type Value struct {
	Pos          lexer.Position
	Deref        *string       `Asterisk @Ident`
	Ref          *string       `| Ref @Ident`
	FunctionCall *FunctionCall `| @@`
	Literal      *Literal      `| @@`
}
type Literal struct {
	Pos    lexer.Position
	Number *float64 `@Number`
	Ident  *string  `| @Ident`
	String *string  `| @String`
	Bool   *bool    `| ( @True | @False )`
	Nil    bool     `| @Nil` // NOTE: might be wrong
}

var luminaLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Ellipsis", `\.\.\.`},
	{"Fn", `fn`},
	{"Extern", `extern`},
	{"Static", `static`},
	{"Const", `const`},
	{"Var", `var`},
	{"Return", `return`},
	{"Number", `(\d*\.)?\d+`},
	{"Dot", `\.`},
	{"String", `\"(?:[^\"]|\\.)*\"`},
	{"Whitespace", `[ \t\n]+`},
	{"LParen", `\(`},
	{"RParen", `\)`},
	{"LBracket", `\[`},
	{"RBracket", `\]`},
	{"LBrace", `{`},
	{"RBrace", `}`},
	{"Colon", `:`},
	{"Comma", `,`},
	{"True", `true`},
	{"False", `false`},
	{"Semicolon", `;`},
	{"Equals", `=`},
	{"Asterisk", `\*`},
	{"FSlash", `/`},
	{"Plus", `\+`},
	{"Minus", `-`},
	{"Ref", `&`},
	{"Nil", `nil`},
	{"Ident", `[a-zA-Z]\w*`}, // NOTE: THIS NEEDS TO GO LATE SO IT DOESNT CONSUME EVERYTHING ELSE!!!!
})
