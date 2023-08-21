@fmt = global [24 x i8] c"sum=%1.1f\0Asecond=%2.2f\0A\00"

declare i32 @printf(i8* %fmt, ...)

define void @main() {
entry:
	%0 = fadd float 2.0, 2.0
	%1 = call i32 (i8*, ...) @printf([24 x i8]* @fmt, float 10.0, float %0)
	ret void
}

define float @getNum(i32 %a) {
entry:
	%0 = fadd float 3.0, 10.0
	%1 = fsub float 1.0, %0
	%2 = fadd float 2.0, %1
	ret float %2
}
