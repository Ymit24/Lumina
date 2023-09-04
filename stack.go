package main

type Stack[T interface{}] struct {
	Top *Node[T]
}

type Node[T interface{}] struct {
	value T
	next  *Node[T]
}

func NewStack[T interface{}]() Stack[T] {
	return Stack[T]{
		Top: nil,
	}
}

func (stack *Stack[T]) Push(value T) {
	next := Node[T]{
		value: value,
		next:  stack.Top,
	}
	stack.Top = &next
}
func (stack *Stack[T]) Pop() T {
	if stack.Top == nil {
		panic("Failed to pop empty stack!")
	}

	top := stack.Top
	stack.Top = top.next

	return top.value
}

func (stack *Stack[T]) Peek() *T {
	return &stack.Top.value
}

func (stack *Stack[T]) PeekNode() *Node[T] {
	return stack.Top
}
