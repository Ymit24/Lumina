package main

import (
	"fmt"
	"math"

	"github.com/llir/llvm/ir/types"
)

var PrimativeTypes = map[string]types.Type{
	"i8":  types.I8,
	"i16": types.I16,
	"i32": types.I32,
	"i64": types.I64,
	"u8":  types.I8,  // NOTE: Signedness is weird in LLVM
	"u16": types.I16, // NOTE: Signedness is weird in LLVM
	"u32": types.I32, // NOTE: Signedness is weird in LLVM
	"u64": types.I64, // NOTE: Signedness is weird in LLVM
	"f32": types.Float,
}

func GetLLVMType(raw Type) (types.Type, error) {
	if raw.Array != nil {
		inner, err := GetLLVMType(raw.Array.Type)
		if err != nil {
			return nil, err
		}
		return types.NewPointer(inner), nil
	}
	var typeName string
	if raw.Inner != nil {
		typeName = raw.Inner.Name
	} else {
		CompileError(raw.Pos, fmt.Errorf("No inner type!"))
	}
	primative, ok := PrimativeTypes[typeName]
	if !ok {
		return nil, fmt.Errorf("Type: `%# v` is not implemented!", raw)
	}
	if raw.Array != nil {
		return types.NewPointer(primative), nil
	}
	return primative, nil
}

func (expr *Expression) GetType() (types.Type, error) {
	leftType, err := expr.Term.GetType()
	if err != nil {
		return nil, err
	}
	if expr.Next != nil {
		rightType, err := expr.Next.GetType()
		if err != nil {
			return nil, err
		}
		if leftType != rightType {
			// TODO: CHECK FOR IMPLICIT CASTING AVAILABILITY (e.g. f32 -> i32)

			if CanCast(leftType, rightType) == false {
				return nil, fmt.Errorf(
					"Expression's left and right sides have different types",
				)
			}
			return TypePrecedence(leftType, rightType), nil
		}
	}
	return leftType, nil
}

func TypePrecedence(src types.Type, dst types.Type) types.Type {
	fmt.Printf("Checking precedence, %# v \t\t%#v", src, dst)
	// TODO: MAKE THIS SO MUCH BETTER
	if types.IsFloat(src) || types.IsFloat(dst) {
		return types.Float
	}
	return src
}

func CanCast(src types.Type, dst types.Type) bool {
	// TODO: IMPLEMENT BETTER WAY TO DO THIS
	if types.IsFloat(src) && types.IsInt(dst) {
		return true
	}
	if types.IsFloat(dst) && types.IsInt(src) {
		return true
	}
	return false
}

func (t *Term) GetType() (types.Type, error) {
	leftType, err := t.Factor.GetType()
	if err != nil {
		return nil, err
	}
	if t.Next != nil {
		rightType, err := t.Next.GetType()
		if err != nil {
			return nil, err
		}
		if leftType != rightType {
			// TODO: CHECK FOR IMPLICIT CASTING AVAILABILITY (e.g. f32 -> i32)
			return nil, fmt.Errorf(
				"Term's left and right sides have different types",
			)
		}
	}
	return leftType, nil
}
func (val *Value) GetType() (types.Type, error) {
	if val.Deref != nil || val.Ref != nil {
		CompileError(
			val.Pos,
			fmt.Errorf("Unimplemented type inference!"),
		)
		// NOTE: SHOULD NOT BE NEEDED DUE TO PANIC
		return nil, fmt.Errorf("Unimplemented type inference!")
	}
	if val.Literal != nil {
		return val.Literal.GetType()
	}
	if val.FunctionCall != nil {
		return val.FunctionCall.GetType()
	}
	return nil, nil
}
func (lit *Literal) GetType() (types.Type, error) {
	if lit.Number != nil {
		if lit.IsFloat() {
			return types.Float, nil
		}
		return types.I32, nil
	}
	if lit.String != nil {
		return types.I8Ptr, nil
	}
	if lit.Bool != nil {
		return types.I1, nil
	}
	if lit.Nil {
		return nil, fmt.Errorf("Can't gain type from nil through inference.")
	}
	return nil, nil
}
func (fc *FunctionCall) GetType() (types.Type, error) {
	return nil, nil
}

func (lit *Literal) IsFloat() bool {
	if lit.Number == nil {
		CompileError(
			lit.Pos,
			fmt.Errorf(
				"Can't check if Literal is a float if it's not a number type!"+
					" Found: %# v", *lit,
			),
		)
	}
	return *lit.Number != math.Trunc(*lit.Number)
}
