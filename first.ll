%Data = type { i32 }
%NewData = type { i32 }

@"fmtfirst.la:8:12" = global [16 x i8] c"A original: %f\0A\00"
@"fmtfirst.la:10:12" = global [15 x i8] c"A changed: %f\0A\00"
@"fmtfirst.la:13:16" = global [7 x i8] c"B: %d\0A\00"
@"fmtfirst.la:15:12" = global [26 x i8] c"A sum=%2.2f\0Asecond=%1.2f\0A\00"

declare i32 @printf(i8* %fmt, ...)

define void @main() {
entry:
	%0 = alloca float
	%1 = sitofp i32 4 to float
	%2 = fadd float %1, 4.5
	store float %2, float* %0
	%3 = load float, float* %0
	%4 = call i32 (i8*, ...) @printf([16 x i8]* @"fmtfirst.la:8:12", float %3)
	%5 = alloca float
	%6 = load float, float* %0
	%7 = sitofp i32 20 to float
	%8 = fadd float %6, %7
	store float %8, float* %5
	%9 = load float, float* %5
	%10 = call i32 (i8*, ...) @printf([15 x i8]* @"fmtfirst.la:10:12", float %9)
	%11 = alloca i32
	store i32 20, i32* %11
	%12 = load i32, i32* %11
	%13 = call i32 (i8*, ...) @printf([7 x i8]* @"fmtfirst.la:13:16", i32 %12)
	%14 = load float, float* %0
	%15 = sitofp i32 2 to float
	%16 = fadd float 0x4001999980000000, %15
	%17 = call i32 (i8*, ...) @printf([26 x i8]* @"fmtfirst.la:15:12", float %14, float %16)
	ret void
}

define i32 @getNum(i32 %a) {
entry:
	%0 = add i32 3, 10
	%1 = sub i32 1, %0
	%2 = add i32 2, %1
	ret i32 %2
}
