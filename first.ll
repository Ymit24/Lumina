%Data = type { i32, float }
%ComplexData = type { %Data }

declare i32 @printf(i8* %fmt, ...)

define void @main() {
entry:
	%0 = alloca %Data
	%1 = alloca %Data
	%2 = getelementptr %Data, %Data* %1, i32 0, i32 0
	store i32 100, i32* %2
	%3 = sitofp i32 3 to float
	%4 = fadd float 200.5, %3
	%5 = getelementptr %Data, %Data* %1, i32 0, i32 1
	store float %4, float* %5
	%6 = load %Data, %Data* %1
	store %Data %6, %Data* %0
	%7 = alloca %ComplexData
	%8 = alloca %ComplexData
	%9 = alloca %Data
	%10 = getelementptr %Data, %Data* %9, i32 0, i32 0
	store i32 2, i32* %10
	%11 = getelementptr %Data, %Data* %9, i32 0, i32 1
	store float 3.5, float* %11
	%12 = load %Data, %Data* %9
	%13 = getelementptr %ComplexData, %ComplexData* %8, i32 0, i32 0
	store %Data %12, %Data* %13
	%14 = load %ComplexData, %ComplexData* %8
	store %ComplexData %14, %ComplexData* %7
	%15 = alloca i32
	%16 = load %Data, %!s(<nil>)
	%17 = bitcast %Data %16 to i32
	store i32 %17, i32* %15
	ret void
}

define %Data @getData() {
entry:
	%0 = alloca %Data
	%1 = getelementptr %Data, %Data* %0, i32 0, i32 0
	store i32 200, i32* %1
	%2 = getelementptr %Data, %Data* %0, i32 0, i32 1
	store float 3.5, float* %2
	%3 = load %Data, %Data* %0
	ret %Data %3
}
