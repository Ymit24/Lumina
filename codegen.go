package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/golang-collections/collections/stack"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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

type ExpressionStep struct {
	Identifier string     // The %2
	Type       types.Type // the 'float', 'i32, etc..
}

func (term *Term) Visit() value.Value {
	// TODO: code gen factor independently
	if term.Factor.Literal != nil {
		lit := *term.Factor.Literal
		if lit.Number != nil {
			return constant.NewFloat(types.Float, *lit.Number)
		}
		CompileError(
			term.Factor.Literal.Pos,
			fmt.Errorf("Unimplemented literal. Found: %# v", lit),
		)
	}
	CompileError(
		term.Factor.Pos,
		fmt.Errorf("Unimplemented factor."),
	)
	return nil
}

func (expr *Expression) Visit() value.Value {
	cBlock := blockStack.Peek().(*ir.Block)
	left := expr.Term.Visit()
	if expr.AddSub != nil {
		// NOTE: For NOW, addition is between ints or floats
		if expr.Next == nil {
			CompileError(
				expr.Pos,
				fmt.Errorf("No next expression in add/sub expression!"),
			)
		}
		right := expr.Next.Visit()
		if *expr.AddSub == "+" {
			return cBlock.NewAdd(
				left,
				right,
			)
		} else if *expr.AddSub == "-" {
			return cBlock.NewSub(
				left,
				right,
			)
		} else {
			CompileError(
				expr.Pos,
				fmt.Errorf("Unrecognized addsub operator: %s", *expr.AddSub),
			)
		}
	} else {
		return left
	}
	CompileError(
		expr.Pos,
		fmt.Errorf("Unimplemented"),
	)
	return nil
	/*
	   entry:
	   %1 = alloca float
	   %2 = fadd float 3.14, float 2
	   %3 = fdiv float 1, float 1
	   %4 = fmul float 5, float %3
	   %5 = fsub 2, %4

	   store float %5, float* %1

	   entry:
	   %1 = alloca float
	   store float 0.14, float* %1
	*/
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

	if stmt.Expression != nil {
		stmt.Expression.Visit()
	}

	cBlock.NewRet(nil)
}

func CompileError(pos lexer.Position, err error) {
	panic(fmt.Sprintf("Failed to compile! At %s reached error %s", pos, err.Error()))
}
