package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/golang-collections/collections/stack"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

var mod *ir.Module
var functionStack stack.Stack
var blockStack stack.Stack

func (p *Program) Visit() error {
	mod = ir.NewModule()
	for _, stmt := range p.Statements {
		stmt.Visit()
	}
	fmt.Printf("Module:\n%s\n", mod.String())
	return nil
}

func (stmt *HighLevelStatement) Visit() {
	if stmt.Extern != nil {
		stmt.Extern.Visit()
	}
	if stmt.Function != nil {
		stmt.Function.Visit()
	}
}

func (stmt *Extern) Visit() {
	fmt.Printf("Found extern: %s\n", stmt.Signiture.Name)
	stmt.Signiture.Visit()
}

func (sig *FunctionSignature) Visit() *ir.Func {
	name := sig.Name
	var returnType types.Type
	var err error
	if sig.ReturnType == nil {
		returnType = types.Void
	} else {
		returnType, err = GetLLVMType(*sig.ReturnType)
		if err != nil {
			CompileError(sig.ReturnType.Pos, err)
		}
	}
	var params []*ir.Param
	for _, arg := range sig.Args {
		argType, err := GetLLVMType(arg.Type)
		if err != nil {
			CompileError(arg.Pos, err)
		}
		params = append(
			params,
			ir.NewParam(arg.Name, argType),
		)
	}
	return mod.NewFunc(
		name,
		returnType,
		params...,
	)
}

func (stmt *Function) Visit() {
	fmt.Printf("Found function: %s\n", stmt.Signature.Name)
	fn := stmt.Signature.Visit()

	functionStack.Push(fn)
	stmt.Block.Visit("entry")
	functionStack.Pop()
}

func (blk *CodeBlock) Visit(name string) {
	fmt.Printf("Found code block.\n")
	cFunc := functionStack.Peek().(*ir.Func)

	block := cFunc.NewBlock(name)
	blockStack.Push(block)

	for _, stmt := range blk.Statements {
		stmt.Visit()
	}

	blockStack.Pop()
}

func (stmt *Statement) Visit() {
	if stmt.VariableAssignment != nil {
		stmt.VariableAssignment.Visit()
	} else if stmt.FunctionCall != nil {
		stmt.FunctionCall.Visit()
	} else if stmt.Return != nil {
		stmt.Return.Visit()
	}
}

func (expr *Expression) Visit() {
}

func (stmt *FunctionCall) Visit() {
	fmt.Printf("Found function call for function: %s\n", stmt.FunctionName)
}
func (stmt *VariableAssignment) Visit() {
	fmt.Printf(
		"Found variable assignment of mutability %s and name %s\n",
		stmt.Mutability,
		stmt.Name,
	)
}
func (stmt *Return) Visit() {
	fmt.Printf("Found return.\n")
	cBlock := blockStack.Peek().(*ir.Block)

	stmt.Expression.Visit()

	cBlock.NewRet(nil)
}

func CompileError(pos lexer.Position, err error) {
	panic(fmt.Sprintf("Failed to compile! At %s reached error %s", pos, err.Error()))
}
