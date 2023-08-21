package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/golang-collections/collections/stack"
	"github.com/kr/pretty"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type CodeGenerator struct {
	module        *ir.Module
	functionStack stack.Stack
	blockStack    stack.Stack
	scopeStack    stack.Stack
}

func (c *CodeGenerator) currentFunction() *ir.Func {
	return c.functionStack.Peek().(*ir.Func)
}

func (c *CodeGenerator) currentBlock() *ir.Block {
	return c.blockStack.Peek().(*ir.Block)
}

func (c *CodeGenerator) currentScope() Scope {
	return c.scopeStack.Peek().(Scope)
}

var functionDeclarations map[string]*ir.Func

func (c *CodeGenerator) VisitProgram(p *Program) string {
	c.module = ir.NewModule()

	functionDeclarations = make(map[string]*ir.Func)

	for _, stmt := range p.Statements {
		c.VisitHighLevelStatement(stmt)
	}

	fmt.Printf("Module:\n%s\n", c.module.String())

	return c.module.String()
}

func (c *CodeGenerator) VisitHighLevelStatement(stmt *HighLevelStatement) {
	if stmt.Extern != nil {
		c.VisitExtern(stmt.Extern)
	}
	if stmt.Function != nil {
		c.VisitFunction(stmt.Function)
	}
}

func (c *CodeGenerator) VisitExtern(stmt *Extern) {
	fmt.Printf("Found extern: %s\n", stmt.Signiture.Name)
	c.VisitFunctionSignature(&stmt.Signiture)
}

func (c *CodeGenerator) VisitFunctionSignature(sig *FunctionSignature) *ir.Func {
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
	isVariadic := false
	for _, arg := range sig.Args {
		argType, err := GetLLVMType(arg.Type)
		if err != nil {
			CompileError(arg.Pos, err)
		}
		fmt.Printf("Arg: %# v\n", pretty.Formatter(arg))
		params = append(
			params,
			ir.NewParam(arg.Name, argType),
		)
		if arg.Type.Array != nil && arg.Type.Array.IsSpread {
			isVariadic = true
		}
	}
	funcDecl := c.module.NewFunc(
		name,
		returnType,
		params...,
	)
	funcDecl.Sig.Variadic = isVariadic

	functionDeclarations[name] = funcDecl
	return funcDecl
}

func (c *CodeGenerator) VisitFunction(stmt *Function) {
	fmt.Printf(
		"Found function: %s, %# v\n",
		stmt.Signature.Name,
		pretty.Formatter(stmt.Signature.Args),
	)
	fn := c.VisitFunctionSignature(&stmt.Signature)

	c.functionStack.Push(fn)
	c.scopeStack.Push(Scope{
		Variables: make(map[string]Variable),
	})
	c.VisitBlock(&stmt.Block, "entry")
	c.scopeStack.Pop()
	c.functionStack.Pop()
}

func (c *CodeGenerator) VisitBlock(blk *CodeBlock, name string) {
	fmt.Printf("Found code block.\n")
	cFunc := c.functionStack.Peek().(*ir.Func)

	block := cFunc.NewBlock(name)
	c.blockStack.Push(block)

	for _, stmt := range blk.Statements {
		c.VisitStatement(&stmt)
	}

	c.blockStack.Pop()
}

func (c *CodeGenerator) VisitStatement(stmt *Statement) {
	if stmt.VariableAssignment != nil {
		c.VisitVariableAssignment(stmt.VariableAssignment)
	} else if stmt.FunctionCall != nil {
		c.VisitFunctionCall(stmt.FunctionCall)
	} else if stmt.Return != nil {
		c.VisitReturn(stmt.Return)
	}
}

type ExpressionStep struct {
	Identifier string     // The %2
	Type       types.Type // the 'float', 'i32, etc..
}

func fixString(raw string) string {
	t := strings.ReplaceAll(raw, `\n`, "\n")
	t = strings.ReplaceAll(t, `"`, "")
	return t
}

func (c *CodeGenerator) VisitTerm(term *Term) value.Value {
	// TODO: code gen factor independently
	if term.Factor.Literal != nil {
		lit := *term.Factor.Literal
		if lit.Number != nil {
			if *lit.Number == math.Trunc(*lit.Number) {
				return constant.NewInt(types.I32, int64(*lit.Number))
			}
			return constant.NewFloat(types.Float, *lit.Number)
		} else if lit.String != nil {
			fmt.Printf("Found string constant: `%s`\n", *lit.String)
			constValue := constant.NewCharArrayFromString(fixString(*lit.String) + "\x00")
			gblDef := c.module.NewGlobalDef("fmt", constValue)
			return gblDef
		} else if lit.Ident != nil {
			// TODO: REPLACE THIS WITH ACTUAL LOOK UP
			return constant.NewCharArrayFromString(*lit.Ident)
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

func (c *CodeGenerator) VisitExpression(expr *Expression) value.Value {
	cBlock := c.blockStack.Peek().(*ir.Block)
	left := c.VisitTerm(&expr.Term)
	if expr.AddSub != nil {
		// NOTE: For NOW, addition is between ints or floats
		if expr.Next == nil {
			CompileError(
				expr.Pos,
				fmt.Errorf("No next expression in add/sub expression!"),
			)
		}
		right := c.VisitExpression(expr.Next)
		if *expr.AddSub == "+" {
			leftIsFloat := types.IsFloat(left.Type())
			rightIsFloat := types.IsFloat(right.Type())
			fmt.Printf("Is Float: %t, %t\n", leftIsFloat, rightIsFloat)
			if leftIsFloat || rightIsFloat {
				if !leftIsFloat {
					left = cBlock.NewBitCast(left, types.Float)
				}
				if !rightIsFloat {
					right = cBlock.NewBitCast(right, types.Float)
				}
				return cBlock.NewFAdd(
					left,
					right,
				)
			} else {
				return cBlock.NewAdd(
					left,
					right,
				)
			}
		} else if *expr.AddSub == "-" {
			leftIsFloat := types.IsFloat(left.Type())
			rightIsFloat := types.IsFloat(right.Type())
			if leftIsFloat || rightIsFloat {
				return cBlock.NewFSub(
					left,
					right,
				)
			} else {
				return cBlock.NewSub(
					left,
					right,
				)
			}
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
}

func (c *CodeGenerator) VisitFunctionCall(stmt *FunctionCall) value.Value {
	fmt.Printf("Found function call for function: %s\n", stmt.FunctionName)
	cBlock := c.currentBlock()

	funcDecl, ok := functionDeclarations[stmt.FunctionName]
	if !ok {
		CompileError(
			stmt.Pos,
			fmt.Errorf("Unknown function %s!", stmt.FunctionName),
		)
	}

	if len(stmt.Args) != 0 {
		var argExprs []value.Value
		for _, arg := range stmt.Args {
			argExprs = append(argExprs, c.VisitExpression(&arg))
		}
		return cBlock.NewCall(
			funcDecl,
			argExprs...,
		)
	}
	return cBlock.NewCall(
		funcDecl,
	)
}
func (c *CodeGenerator) VisitVariableAssignment(stmt *VariableAssignment) {
	fmt.Printf(
		"Found variable assignment of mutability %s and name %s\n",
		stmt.Mutability,
		stmt.Name,
	)

	fmt.Printf("%# v\n", stmt)

	vType, err := GetLLVMType(*stmt.Type)
	if err != nil {
		CompileError(stmt.Pos, err)
	}

	if _, ok := c.currentScope().Variables[stmt.Name]; ok {
		CompileError(
			stmt.Pos,
			fmt.Errorf("Variable %s already exists in scope!", stmt.Name),
		)
	}

	alloc := c.currentBlock().NewAlloca(vType)
	variable := Variable{
		Name:       stmt.Name,
		Mutability: stmt.Mutability,
		Type:       vType,
		Address:    alloc,
	}

	switch stmt.Mutability {
	case "static":
		CompileError(
			stmt.Pos,
			fmt.Errorf("Statics are not implemented"),
		)
	case "const":
		c.currentScope().Variables[stmt.Name] = variable
	case "var":
		c.currentScope().Variables[stmt.Name] = variable
	}

	exprResult := c.VisitExpression(&stmt.Expression)
	c.currentBlock().NewStore(exprResult, alloc)
}
func (c *CodeGenerator) VisitReturn(stmt *Return) {
	fmt.Printf("Found return.\n")
	cBlock := c.currentBlock()

	if stmt.Expression != nil {
		val := c.VisitExpression(stmt.Expression)
		cBlock.NewRet(val)
	} else {
		cBlock.NewRet(nil)
	}
}

func CompileError(pos lexer.Position, err error) {
	panic(fmt.Sprintf("Failed to compile! At %s reached error %s", pos, err.Error()))
}
