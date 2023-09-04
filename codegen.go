package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/kr/pretty"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type CodeGenerator struct {
	module                   *ir.Module
	rootFunctionDeclarations map[string]*ir.Func
	scopeStack               Stack[Scope]
}

func (c *CodeGenerator) currentBlock() *ir.Block {
	return c.currentScope().GeneratingBlock
}

func (c *CodeGenerator) currentScope() *Scope {
	return c.scopeStack.Peek()
}

func NewCodeGenerator() CodeGenerator {
	scopeStack := NewStack[Scope]()
	scopeStack.Push(Scope{
		Type:             Global,
		Variables:        make(map[string]Variable),
		StructDefintions: make(map[string]Struct),
	})
	return CodeGenerator{
		scopeStack:               scopeStack,
		rootFunctionDeclarations: make(map[string]*ir.Func),
	}
}

func (c *CodeGenerator) VisitProgram(p *Program) string {
	c.module = ir.NewModule()

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
	if stmt.Struct != nil {
		c.VisitStructDefinition(stmt.Struct)
	}
}

func (c *CodeGenerator) VisitStructDefinition(structDefinition *StructDefinition) {
	llvmTypes, err := c.GetStructLLVMTypes(structDefinition.Body.Fields)
	if err != nil {
		CompileError(
			structDefinition.Pos,
			fmt.Errorf(
				"Failed to get llvm types for struct definition '%s'. Error: %s",
				structDefinition.Name,
				err.Error(),
			),
		)
	}
	structType := types.NewStruct(llvmTypes...)
	// TODO: MAKE SURE THE NAME IS ALWAYS UNIQUE. THIS NEEDS TO HANDLE
	// TOP LEVEL STRUCTS, LOCAL STRUCTS, ANON STRUCTS, GENERICS, ETC...
	structTypeDef := c.module.NewTypeDef(structDefinition.Name, structType)

	var structFields []string
	for _, field := range structDefinition.Body.Fields {
		structFields = append(structFields, field.Name)
	}

	c.currentScope().StructDefintions[structDefinition.Name] = Struct{
		Name:          structDefinition.Name,
		TypeDef:       structTypeDef,
		StructDef:     structType,
		OrderedFields: structFields,
	}
}

func (c *CodeGenerator) NewScope(scopeType ScopeType, generatingBlock *ir.Block) {
	var generatingBlockToUse *ir.Block
	if generatingBlock == nil {
		generatingBlockToUse = c.currentScope().GeneratingBlock
	} else {
		generatingBlockToUse = generatingBlock
	}
	c.scopeStack.Push(Scope{
		Type:             scopeType,
		Variables:        make(map[string]Variable),
		StructDefintions: make(map[string]Struct),
		GeneratingBlock:  generatingBlockToUse,
	})
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
		returnType, err = c.GetLLVMType(*sig.ReturnType)
		if err != nil {
			CompileError(sig.ReturnType.Pos, err)
		}
	}
	var params []*ir.Param
	isVariadic := false
	for _, arg := range sig.Args {
		argType, err := c.GetLLVMType(arg.Type)
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

	c.rootFunctionDeclarations[name] = funcDecl
	return funcDecl
}

func (c *CodeGenerator) VisitFunction(stmt *Function) {
	fmt.Printf(
		"Found function: %s, %# v\n",
		stmt.Signature.Name,
		pretty.Formatter(stmt.Signature.Args),
	)
	fn := c.VisitFunctionSignature(&stmt.Signature)
	functionBlock := fn.NewBlock("entry")

	c.NewScope(FunctionScope, functionBlock)

	c.VisitBlock(&stmt.Block)
	c.scopeStack.Pop()
}

// NOTE: In GENERAL, code blocks do not actually correspond to new llvm blocks.
func (c *CodeGenerator) VisitBlock(blk *CodeBlock) {
	for _, stmt := range blk.Statements {
		c.VisitStatement(&stmt)
	}
}

func (c *CodeGenerator) VisitStatement(stmt *Statement) {
	if stmt.VariableDeclaration != nil {
		c.VisitVariableDeclaration(stmt.VariableDeclaration)
	} else if stmt.VariableAssignment != nil {
		c.VisitVariableAssignment(stmt.VariableAssignment)
	} else if stmt.FunctionCall != nil {
		c.VisitFunctionCall(stmt.FunctionCall)
	} else if stmt.Return != nil {
		c.VisitReturn(stmt.Return)
	} else if stmt.ScopeBlock != nil {
		c.NewScope(CodeBlockScope, nil)
		c.scopeStack.Push(Scope{
			Type:      CodeBlockScope,
			Variables: make(map[string]Variable),
			// Scope blocks inherit their generating block
			GeneratingBlock: c.currentScope().GeneratingBlock,
		})
		c.VisitBlock(stmt.ScopeBlock)
		c.scopeStack.Pop()
	} else if stmt.StructDefinition != nil {
		c.VisitStructDefinition(stmt.StructDefinition)
	} else {
		CompileError(
			stmt.Pos,
			fmt.Errorf("Unknown statement type."),
		)
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
			if lit.IsFloat() {
				return constant.NewFloat(types.Float, *lit.Number)
			}
			return constant.NewInt(types.I32, int64(*lit.Number))
		} else if lit.String != nil {
			constValue := constant.NewCharArrayFromString(
				fixString(*lit.String) + "\x00",
			)
			gblDef := c.module.NewGlobalDef("fmt"+lit.Pos.String(), constValue)
			return gblDef
		} else if lit.Ident != nil {
			variable, err := c.getVariable(*lit.Ident)
			if err != nil {
				CompileError(
					lit.Pos,
					err,
				)
			}
			return c.currentScope().GeneratingBlock.NewLoad(variable.Type, variable.Address)
		} else if lit.Struct != nil {
			fmt.Printf("Found struct literal %# v\n", pretty.Formatter(lit.Struct))
			structType, err := c.getStructDefinition(lit.Struct.Name)
			if err != nil {
				CompileError(
					lit.Struct.Pos,
					fmt.Errorf(
						"Failed to instantiate Struct literal. Error: %s\n",
						err.Error(),
					),
				)
			}
			fmt.Printf("Found the type of the struct. Creating.\n")

			cBlock := c.currentBlock()
			structAddress := cBlock.NewAlloca(structType.TypeDef)
			// NOTE: THIS ASSUMES THAT THE FIELDS ARE IN THE ORDER THEY'RE DEFINED
			for _, field := range lit.Struct.Body.Fields {
				currentExpression := c.VisitExpression(&field.Expression)
				index, err := structType.getFieldIndex(field.Name)
				if err != nil {
					CompileError(
						field.Pos,
						fmt.Errorf(
							"Failed to get field index. Error: %s\n",
							err.Error(),
						),
					)
				}
				currentFieldAddress := cBlock.NewGetElementPtr(
					structType.TypeDef,
					structAddress,
					constant.NewIndex(constant.NewInt(types.I32, 0)),
					constant.NewIndex(constant.NewInt(types.I32, index)),
				)
				// castedElementAddress := cBlock.NewBitCast(currentFieldAddress, types.NewPointer(structType.StructDef.Fields[index]))
				cBlock.NewStore(currentExpression, currentFieldAddress)
			}
			return cBlock.NewLoad(structType.TypeDef, structAddress)
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
	cBlock := c.currentScope().GeneratingBlock
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
			if leftIsFloat || rightIsFloat {
				if !leftIsFloat {
					left = cBlock.NewSIToFP(left, types.Float)
					// left = cBlock.NewBitCast(left, types.Float)
				}
				if !rightIsFloat {
					right = cBlock.NewSIToFP(right, types.Float)
					// right = cBlock.NewBitCast(right, types.Float)
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
				if !leftIsFloat {
					left = cBlock.NewSIToFP(left, types.Float)
					// left = cBlock.NewBitCast(left, types.Float)
				}
				if !rightIsFloat {
					right = cBlock.NewSIToFP(right, types.Float)
					// right = cBlock.NewBitCast(right, types.Float)
				}
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
	cBlock := c.currentScope().GeneratingBlock

	funcDecl, ok := c.rootFunctionDeclarations[stmt.FunctionName]
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

func (c *CodeGenerator) getVariable(name string) (Variable, error) {
	scopeNode := c.scopeStack.PeekNode()
	for scopeNode != nil {
		variables := scopeNode.value.Variables
		if _, ok := variables[name]; ok {
			fmt.Printf("Found variable!\n")
			return variables[name], nil
		}
		fmt.Printf("Didn't find variable in current scope.\n")
		scopeNode = scopeNode.next
	}
	return Variable{}, fmt.Errorf("No variable found by that name in any valid scope!")
}

func (c *CodeGenerator) getStructDefinition(name string) (Struct, error) {
	scopeNode := c.scopeStack.PeekNode()
	for scopeNode != nil {
		structDefinitions := scopeNode.value.StructDefintions
		if _, ok := structDefinitions[name]; ok {
			fmt.Printf("Found struct definition!\n")
			return structDefinitions[name], nil
		}
		fmt.Printf("Didn't find struct definition in current scope.\n")
		scopeNode = scopeNode.next
	}
	return Struct{}, fmt.Errorf("No struct definition found by that name in any valid scope!")
}

func (c *CodeGenerator) VisitVariableAssignment(stmt *VariableAssignment) {
	cBlock := c.currentBlock()
	variable, err := c.getVariable(stmt.Name)
	if err != nil {
		CompileError(
			stmt.Pos,
			fmt.Errorf("Trying to assign to undeclared variable!"),
		)
	}
	value := c.VisitExpression(&stmt.Expression)
	cBlock.NewStore(value, variable.Address)
}

func (c *CodeGenerator) VisitVariableDeclaration(stmt *VariableDeclaration) {
	fmt.Printf(
		"Found variable declaration of mutability %s and name %s\n",
		stmt.Mutability,
		stmt.Name,
	)

	var llvmType types.Type
	var err error

	fmt.Printf("%# v\n", stmt)
	if stmt.Type == nil {
		// NOTE: Try type inference
		llvmType, err = stmt.Expression.GetType()
		if err != nil {
			CompileError(
				stmt.Pos,
				fmt.Errorf(
					"Unable to infer type of expression for "+
						"variable assignment. Error: %s",
					err.Error(),
				),
			)
		}
		fmt.Printf(
			"\nType Inference returned: %# v for %# v\n\n",
			llvmType,
			pretty.Formatter(stmt.Expression),
		)
	} else {
		llvmType, err = c.GetLLVMType(*stmt.Type)
	}
	if err != nil {
		CompileError(stmt.Pos, err)
	}

	if _, ok := c.currentScope().Variables[stmt.Name]; ok {
		CompileError(
			stmt.Pos,
			fmt.Errorf("Variable %s already exists in scope!", stmt.Name),
		)
	}

	cBlock := c.currentScope().GeneratingBlock
	alloc := cBlock.NewAlloca(llvmType)
	variable := Variable{
		Name:       stmt.Name,
		Mutability: stmt.Mutability,
		Type:       llvmType,
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
	if exprResult.Type() != alloc.ElemType {
		exprResult = cBlock.NewBitCast(exprResult, alloc.ElemType)
	}
	cBlock.NewStore(exprResult, alloc)
}

func (c *CodeGenerator) VisitReturn(stmt *Return) {
	fmt.Printf("Found return.\n")
	cBlock := c.currentScope().GeneratingBlock

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
