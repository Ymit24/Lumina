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
	store i32 300, i32* %10
	%11 = fsub float 0x40590CCCC0000000, 5.5
	%12 = getelementptr %Data, %Data* %9, i32 0, i32 1
	store float %11, float* %12
	%13 = load %Data, %Data* %9
	%14 = getelementptr %ComplexData, %ComplexData* %8, i32 0, i32 0
	store %Data %13, %Data* %14
	%15 = load %ComplexData, %ComplexData* %8
	store %ComplexData %15, %ComplexData* %7
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
