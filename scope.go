package main

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Scope struct {
	Variables map[string]Variable
}

type Variable struct {
	Name       string
	Mutability string
	Type       types.Type
	Address    value.Value
}
