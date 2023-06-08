package main

import "fmt"

var depth int = 0

func push() {
	fmt.Println("	push %rax")
	depth++
}

func pop(reg string) {
	fmt.Printf("	pop %s\n", reg)
	depth--
}

func genAddr(node *Node) {
	if node.kind != ND_VAR {
		Errorf("error genAddr")
	}

	offset := (int(node.name[0]) - 'a' + 1) * 8
	fmt.Printf("	lea %d(%%rbp), %%rax\n", -offset)
}

func genExpr(node *Node) {
	switch node.kind {
	case ND_NUM:
		fmt.Printf("	mov $%d, %%rax\n", node.val)
		return
	case ND_NEG:
		genExpr(node.lhs)
		fmt.Println("	neg %rax")
		return
	case ND_ASSIGN:
		genAddr(node.lhs)
		push()
		genExpr(node.rhs)
		pop("%rdi")
		fmt.Println("	mov %rax, (%rdi)")
		return
	case ND_VAR:
		genAddr(node)
		fmt.Println("	mov (%rax), %rax")
		return
	}

	genExpr(node.rhs)
	push()
	genExpr(node.lhs)
	pop("%rdi")

	switch node.kind {
	case ND_ADD:
		fmt.Println("	add %rdi, %rax")
		return
	case ND_SUB:
		fmt.Println("	sub %rdi, %rax")
		return
	case ND_MUL:
		fmt.Println("	imul %rdi, %rax")
		return
	case ND_DIV:
		fmt.Println("	cqo")
		fmt.Println("	idiv %rdi")
		return
	case ND_EQ, ND_NE, ND_LT, ND_LE:
		fmt.Println("	cmp %rdi, %rax")

		if node.kind == ND_EQ {
			fmt.Println("	sete %al")
		} else if node.kind == ND_NE {
			fmt.Println("	setne %al")
		} else if node.kind == ND_LT {
			fmt.Println("	setl %al")
		} else {
			fmt.Println("	setle %al")
		}

		fmt.Println("	movzb %al, %rax")
		return
	}

	Errorf("error genExpr")
}

func Codegen(node *Node) {
	fmt.Println(".globl main")
	fmt.Println("main:")

	fmt.Println("	push %rbp")
	fmt.Println("	mov %rsp, %rbp")
	fmt.Println("	sub $208, %rsp")

	for cur := node; cur != nil; cur = cur.next {
		if cur.kind == ND_EXPR_STMT {
			genExpr(cur.lhs)
			if depth != 0 {
				Errorf("error depth")
			}
			continue
		}
		Errorf("error codegen")
	}

	fmt.Println("	mov %rbp, %rsp")
	pop("%rbp")
	fmt.Println("	ret")
}
