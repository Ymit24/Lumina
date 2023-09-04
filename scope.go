package main

import (
	"fmt"

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
	StructDefintions map[string]Struct
	// This is the llvm block to use
	GeneratingBlock *ir.Block
}

type Variable struct {
	Name       string
	Mutability string
	Type       types.Type
	Address    value.Value
}

type Struct struct {
	Name          string
	TypeDef       types.Type
	StructDef     *types.StructType
	OrderedFields []string
}

func (s *Struct) getFieldIndex(name string) (int64, error) {
	for index, field := range s.OrderedFields {
		if field == name {
			fmt.Printf("Found index of %s at %d\n.", name, index)
			return int64(index), nil
		}
	}
	return -1, fmt.Errorf(
		"Failed to find index of field named %s for struct %s.",
		name,
		s.Name,
	)
}
