%Data = type { i32, float }
%ComplexData = type { %Data }
%Epic = type { %Data, %ComplexData, %ComplexData }

@"fmtfirst.la:19:12" = global [13 x i8] c"Data: %d %f\0A\00"
@"fmtfirst.la:34:12" = global [11 x i8] c"Value: %f\0A\00"
@"fmtfirst.la:36:12" = global [11 x i8] c"Value: %f\0A\00"

declare i32 @printf(i8* %fmt, ...)

define %Data @getData() {
entry:
	%0 = alloca %Data
	%1 = getelementptr %Data, %Data* %0, i32 0, i32 0
	store i32 200, i32* %1
	%2 = getelementptr %Data, %Data* %0, i32 0, i32 1
	store float 100.5, float* %2
	%3 = load %Data, %Data* %0
	ret %Data %3
}

define void @main() {
entry:
	%0 = alloca %Data
	%1 = call %Data @getData()
	store %Data %1, %Data* %0
	%2 = getelementptr %Data, %Data* %0, i32 0, i32 0
	%3 = load i32, i32* %2
	%4 = getelementptr %Data, %Data* %0, i32 0, i32 1
	%5 = load float, float* %4
	%6 = call i32 (i8*, ...) @printf([13 x i8]* @"fmtfirst.la:19:12", i32 %3, float %5)
	%7 = alloca %ComplexData
	%8 = alloca %ComplexData
	%9 = alloca %Data
	%10 = getelementptr %Data, %Data* %9, i32 0, i32 0
	store i32 2, i32* %10
	%11 = getelementptr %Data, %Data* %9, i32 0, i32 1
	store float 4.5, float* %11
	%12 = load %Data, %Data* %9
	%13 = getelementptr %ComplexData, %ComplexData* %8, i32 0, i32 0
	store %Data %12, %Data* %13
	%14 = load %ComplexData, %ComplexData* %8
	store %ComplexData %14, %ComplexData* %7
	%15 = alloca %Epic
	%16 = alloca %Epic
	%17 = alloca %Data
	%18 = getelementptr %Data, %Data* %17, i32 0, i32 0
	store i32 20, i32* %18
	%19 = getelementptr %Data, %Data* %17, i32 0, i32 1
	store float 10.25, float* %19
	%20 = load %Data, %Data* %17
	%21 = getelementptr %Epic, %Epic* %16, i32 0, i32 0
	store %Data %20, %Data* %21
	%22 = alloca %ComplexData
	%23 = alloca %Data
	%24 = getelementptr %Data, %Data* %23, i32 0, i32 0
	store i32 0, i32* %24
	%25 = sitofp i32 2 to float
	%26 = fadd float 10.25, %25
	%27 = getelementptr %Data, %Data* %23, i32 0, i32 1
	store float %26, float* %27
	%28 = load %Data, %Data* %23
	%29 = getelementptr %ComplexData, %ComplexData* %22, i32 0, i32 0
	store %Data %28, %Data* %29
	%30 = load %ComplexData, %ComplexData* %22
	%31 = getelementptr %Epic, %Epic* %16, i32 0, i32 1
	store %ComplexData %30, %ComplexData* %31
	%32 = load %ComplexData, %ComplexData* %7
	%33 = getelementptr %Epic, %Epic* %16, i32 0, i32 2
	store %ComplexData %32, %ComplexData* %33
	%34 = load %Epic, %Epic* %16
	store %Epic %34, %Epic* %15
	%35 = getelementptr %ComplexData, %ComplexData* %7, i32 0, i32 0, i32 1
	%36 = load float, float* %35
	%37 = call i32 (i8*, ...) @printf([11 x i8]* @"fmtfirst.la:34:12", float %36)
	%38 = getelementptr %ComplexData, %ComplexData* %7, i32 0, i32 0, i32 1
	store float 10.5, float* %38
	%39 = getelementptr %ComplexData, %ComplexData* %7, i32 0, i32 0, i32 1
	%40 = load float, float* %39
	%41 = call i32 (i8*, ...) @printf([11 x i8]* @"fmtfirst.la:36:12", float %40)
	ret void
}
