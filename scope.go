package main

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type ScopeType int

const (
	Global ScopeType = iota
	FunctionScope
	ControlScope   // Catch all for if, for, case, etc..
	CodeBlockScope // e.g. anon blocks
)

func (s ScopeType) String() string {
	return [...]string{"Global", "FunctionScope", "ControlScope", "CodeBlockScope"}[s-1]
}

type Scope struct {
	Type             ScopeType
	Variables        map[string]Variable
	StructDefintions map[string]types.Type
	// This is the llvm block to use
	GeneratingBlock *ir.Block
}

type Variable struct {
	Name       string
	Mutability string
	Type       types.Type
	Address    value.Value
}
