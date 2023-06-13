package codegen

import (
	"3ccgo6/header"
	"fmt"
)

var depth int = 0

func push() {
	fmt.Println("	push %rax")
	depth++
}

func pop(msg string) {
	fmt.Printf("	pop %s\n", msg)
	depth--
}

func genAddr(node *header.Node) {
	if node.Kind != header.NdVar {
		header.Errorf("error genAddr")
	}
	offset := (int(rune(node.Name[0])) - 'a' + 1) * 8
	fmt.Printf("	lea %d(%%rbp), %%rax\n", -offset)
}

func genExpr(node *header.Node) {
	switch node.Kind {
	case header.NdNum:
		fmt.Printf("	mov $%d, %%rax\n", node.Val)
		return
	case header.NdNeg:
		genExpr(node.Lhs)
		fmt.Println("	neg %rax")
		return
	case header.NdAssign:
		genAddr(node.Lhs)
		push()
		genExpr(node.Rhs)
		pop("%rdi")
		fmt.Println("	mov %rax, (%rdi)")
		return
	case header.NdVar:
		genAddr(node)
		fmt.Println("	mov (%rax), %rax")
		return
	}

	genExpr(node.Rhs)
	push()
	genExpr(node.Lhs)
	pop("%rdi")

	switch node.Kind {
	case header.NdAdd:
		fmt.Println("	add %rdi, %rax")
		return
	case header.NdSub:
		fmt.Println("	sub %rdi, %rax")
		return
	case header.NdMul:
		fmt.Println("	imul %rdi, %rax")
		return
	case header.NdDiv:
		fmt.Println("	cqo")
		fmt.Println("	idiv %rdi")
		return
	case header.NdEq, header.NdNe, header.NdLt, header.NdLe:
		fmt.Println("	cmp %rdi, %rax")

		if node.Kind == header.NdEq {
			fmt.Println("	sete %al")
		} else if node.Kind == header.NdNe {
			fmt.Println("	setne %al")
		} else if node.Kind == header.NdLt {
			fmt.Println("	setl %al")
		} else {
			fmt.Println("	setle %al")
		}

		fmt.Println("	movzb %al, %rax")
		return
	}

	header.Errorf("error genexpr")
}

func Codegen(node *header.Node) {
	fmt.Println(".globl main")
	fmt.Println("main:")

	fmt.Println("	push %rbp")
	fmt.Println("	mov %rsp, %rbp")
	fmt.Printf("	sub $208, %%rsp\n")

	for cur := node; cur != nil; cur = cur.Next {
		if cur.Kind != header.NdExprStmt {
			header.Errorf("Error Codegen NdExprStmt")
		}
		genExpr(cur.Lhs)
		if depth != 0 {
			header.Errorf("Error Codegen depth")
		}
	}

	fmt.Println("	mov %rbp, %rsp")
	fmt.Println("	pop %rbp")
	fmt.Println("	ret")
}
