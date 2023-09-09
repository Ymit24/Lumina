%Data = type { i32, float }
%ComplexData = type { %Data }
%Epic = type { %Data, %ComplexData, %ComplexData }

@"fmtfirst.la:23:12" = global [13 x i8] c"Data: %d %f\0A\00"
@"fmtfirst.la:24:12" = global [30 x i8] c"Done something with data: %d\0A\00"
@"fmtfirst.la:39:12" = global [11 x i8] c"Value: %f\0A\00"
@"fmtfirst.la:41:12" = global [11 x i8] c"Value: %f\0A\00"

declare i32 @printf(i8* %fmt, ...)

define %Data @getData() {
entry:
	%0 = alloca %Data
	%1 = getelementptr %Data, %Data* %0, i32 0, i32 0
	store i32 201, i32* %1
	%2 = getelementptr %Data, %Data* %0, i32 0, i32 1
	store float 100.5, float* %2
	%3 = load %Data, %Data* %0
	ret %Data %3
}

define i32 @doSomethingWithData(i32 %i, i32 %x) {
entry:
	%0 = add i32 %i, %x
	ret i32 %0
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
	%6 = call i32 (i8*, ...) @printf([13 x i8]* @"fmtfirst.la:23:12", i32 %3, float %5)
	%7 = getelementptr %Data, %Data* %0, i32 0, i32 0
	%8 = load i32, i32* %7
	%9 = getelementptr %Data, %Data* %0, i32 0, i32 0
	%10 = load i32, i32* %9
	%11 = call i32 @doSomethingWithData(i32 %10, i32 1)
	%12 = call i32 @doSomethingWithData(i32 %8, i32 %11)
	%13 = call i32 (i8*, ...) @printf([30 x i8]* @"fmtfirst.la:24:12", i32 %12)
	%14 = alloca %ComplexData
	%15 = alloca %ComplexData
	%16 = alloca %Data
	%17 = getelementptr %Data, %Data* %16, i32 0, i32 0
	store i32 2, i32* %17
	%18 = getelementptr %Data, %Data* %16, i32 0, i32 1
	store float 4.5, float* %18
	%19 = load %Data, %Data* %16
	%20 = getelementptr %ComplexData, %ComplexData* %15, i32 0, i32 0
	store %Data %19, %Data* %20
	%21 = load %ComplexData, %ComplexData* %15
	store %ComplexData %21, %ComplexData* %14
	%22 = alloca %Epic
	%23 = alloca %Epic
	%24 = alloca %Data
	%25 = getelementptr %Data, %Data* %24, i32 0, i32 0
	store i32 20, i32* %25
	%26 = getelementptr %Data, %Data* %24, i32 0, i32 1
	store float 10.25, float* %26
	%27 = load %Data, %Data* %24
	%28 = getelementptr %Epic, %Epic* %23, i32 0, i32 0
	store %Data %27, %Data* %28
	%29 = alloca %ComplexData
	%30 = alloca %Data
	%31 = getelementptr %Data, %Data* %30, i32 0, i32 0
	store i32 0, i32* %31
	%32 = sitofp i32 2 to float
	%33 = fadd float 10.25, %32
	%34 = getelementptr %Data, %Data* %30, i32 0, i32 1
	store float %33, float* %34
	%35 = load %Data, %Data* %30
	%36 = getelementptr %ComplexData, %ComplexData* %29, i32 0, i32 0
	store %Data %35, %Data* %36
	%37 = load %ComplexData, %ComplexData* %29
	%38 = getelementptr %Epic, %Epic* %23, i32 0, i32 1
	store %ComplexData %37, %ComplexData* %38
	%39 = load %ComplexData, %ComplexData* %14
	%40 = getelementptr %Epic, %Epic* %23, i32 0, i32 2
	store %ComplexData %39, %ComplexData* %40
	%41 = load %Epic, %Epic* %23
	store %Epic %41, %Epic* %22
	%42 = getelementptr %ComplexData, %ComplexData* %14, i32 0, i32 0, i32 1
	%43 = load float, float* %42
	%44 = call i32 (i8*, ...) @printf([11 x i8]* @"fmtfirst.la:39:12", float %43)
	%45 = getelementptr %ComplexData, %ComplexData* %14, i32 0, i32 0, i32 1
	store float 10.5, float* %45
	%46 = getelementptr %ComplexData, %ComplexData* %14, i32 0, i32 0, i32 1
	%47 = load float, float* %46
	%48 = call i32 (i8*, ...) @printf([11 x i8]* @"fmtfirst.la:41:12", float %47)
	ret void
}
