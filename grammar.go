package main

type Program struct {
	Statements []*HighLevelStatement `@@+`
}

type HighLevelStatement struct {
	Extern   *Extern   `@@`
	Function *Function `| @@`
}

type Extern struct {
	Signiture FunctionSignature `"extern" @@ Semicolon`
}

type FunctionSignature struct {
	ReturnType *Type               `@@`
	Name       string              `@Ident`
	Args       []*FunctionArgument `(LParen ( @@ ( Comma @@ )* )? RParen )`
}

type Function struct {
	Signature FunctionSignature `@@`
	Block     CodeBlock         `@@`
}
type Type struct {
	Name             string  `@Ident`
	GenericArguments []*Type `( LParen ( @@ ( Comma @@ )* ) RParen )?`
}
type CodeBlock struct {
	Statements []Statement `LBrace ( @@ )* RBrace`
}
type FunctionArgument struct {
	Name string `@Ident`
	Type Type   `Colon @@`
}
type Statement struct {
	FunctionCall *FunctionCall `@@ Semicolon`
}
type FunctionCall struct {
	FunctionName string  `@Ident`
	Args         []Value `LParen @@ (Comma @@)* RParen`
}
type Value struct {
	Number *float64 ` @Number`
	String *string  `| @String | @Ident`
	Bool   *bool    `| ( @True | @False )`
	Nil    bool     `| @"nil"`
}
