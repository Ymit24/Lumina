@fmt = global [21 x i8] c"sum=%d\0Asecond=%1.2f\0A\00"

declare i32 @printf(i8* %fmt, ...)

define void @main() {
entry:
	%0 = bitcast i32 2 to float
	%1 = fadd float 0x4001999980000000, %0
	%2 = call i32 (i8*, ...) @printf([21 x i8]* @fmt, i32 10, float %1)
	ret void
}

define i32 @getNum(i32 %a) {
entry:
	%0 = add i32 3, 10
	%1 = sub i32 1, %0
	%2 = add i32 2, %1
	ret i32 %2
}
