package main

import (
	"fmt"

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
	typeName := raw.Name
	primative, ok := PrimativeTypes[typeName]
	if ok {
		return primative, nil
	}
	return nil, fmt.Errorf("Type: `%# v` is not implemented!", raw)
}
