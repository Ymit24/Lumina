declare i32 @puts(i8 %str)

define void @main() {
entry:
        %0 = call i32 @puts([16 x i8] c"\22Hello world!\5Cn\22")
        ret void
}
