@fmt = global [24 x i8] c"sum=%2.2f\0Asecond=%1.2f\0A\00"

declare i32 @printf(i8* %fmt, ...)

define void @main() {
entry:
	%0 = alloca float
	%1 = bitcast i32 4 to float
	%2 = fadd float %1, 4.5
	store float %2, float* %0
	%3 = load float, float* %0
	%4 = bitcast i32 2 to float
	%5 = fadd float 0x4001999980000000, %4
	%6 = call i32 (i8*, ...) @printf([24 x i8]* @fmt, float %3, float %5)
	ret void
}

define i32 @getNum(i32 %a) {
entry:
	%0 = add i32 3, 10
	%1 = sub i32 1, %0
	%2 = add i32 2, %1
	ret i32 %2
}
