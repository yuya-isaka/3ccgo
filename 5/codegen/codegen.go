package codegen

import (
	"3ccgo5/header"
	"fmt"
)

var depth int = 0

func push() {
	fmt.Println("	push %rax")
	depth++
}

func pop(reg string) {
	fmt.Printf("	pop %s\n", reg)
	depth--
}

func genAddr(node *header.Node) {
	if node.Kind != header.ND_VAR {
		header.Errorf("error genAddr")
	}

	offset := (int(node.Name[0]) - 'a' + 1) * 8
	fmt.Printf("	lea %d(%%rbp), %%rax\n", -offset)
}

func genExpr(node *header.Node) {
	switch node.Kind {
	case header.ND_NUM:
		fmt.Printf("	mov $%d, %%rax\n", node.Val)
		return
	case header.ND_NEG:
		genExpr(node.Lhs)
		fmt.Println("	neg %rax")
		return
	case header.ND_ASSIGN:
		genAddr(node.Lhs)
		push()
		genExpr(node.Rhs)
		pop("%rdi")
		fmt.Println("	mov %rax, (%rdi)")
		return
	case header.ND_VAR:
		genAddr(node)
		fmt.Println("	mov (%rax), %rax")
		return
	}

	genExpr(node.Rhs)
	push()
	genExpr(node.Lhs)
	pop("%rdi")

	switch node.Kind {
	case header.ND_ADD:
		fmt.Println("	add %rdi, %rax")
		return
	case header.ND_SUB:
		fmt.Println("	sub %rdi, %rax")
		return
	case header.ND_MUL:
		fmt.Println("	imul %rdi, %rax")
		return
	case header.ND_DIV:
		fmt.Println("	cqo")
		fmt.Println("	idiv %rdi")
		return
	case header.ND_EQ, header.ND_NE, header.ND_LT, header.ND_LE:
		fmt.Println("	cmp %rdi, %rax")

		if node.Kind == header.ND_EQ {
			fmt.Println("	sete %al")
		} else if node.Kind == header.ND_NE {
			fmt.Println("	setne %al")
		} else if node.Kind == header.ND_LT {
			fmt.Println("	setl %al")
		} else {
			fmt.Println("	setle %al")
		}

		fmt.Println("	movzb %al, %rax")
		return
	}

	header.Errorf("error genExpr")
}

func Codegen(node *header.Node) {
	fmt.Println(".globl main")
	fmt.Println("main:")

	fmt.Println("	push %rbp")
	fmt.Println("	mov %rsp, %rbp")
	fmt.Println("	sub $208, %rsp")

	for cur := node; cur != nil; cur = cur.Next {
		if cur.Kind == header.ND_EXPR_STMT {
			genExpr(cur.Lhs)
			if depth != 0 {
				header.Errorf("error depth")
			}
			continue
		}
		header.Errorf("error codegen")
	}

	fmt.Println("	mov %rbp, %rsp")
	pop("%rbp")
	fmt.Println("	ret")
}
