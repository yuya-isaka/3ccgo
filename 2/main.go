package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type TokenKind int

const (
	TK_PUNCT TokenKind = iota
	TK_NUM
	TK_EOF
	TK_VAR
)

type Token struct {
	kind TokenKind
	next *Token
	loc  int
	val  int
	str  string
}

func newToken(kind TokenKind, index int) *Token {
	tok := Token{kind: kind, loc: index}
	return &tok
}

var user_input []rune

func error_at(at int, err error) {
	fmt.Fprintf(os.Stderr, "%s\n", string(user_input))
	fmt.Fprintf(os.Stderr, "%*s", at, "")
	fmt.Fprintf(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
	os.Exit(1)
}

func getNum(text []rune, index *int) int {
	result, err := strconv.Atoi(string(text[*index]))
	if err != nil {
		error_at(*index, fmt.Errorf("error getNum"))
	}
	*index++
	for *index < len(text) && unicode.IsDigit(text[*index]) {
		tmp, err := strconv.Atoi(string(text[*index]))
		if err != nil {
			break
		}
		result = result*10 + tmp
		*index++
	}
	return result
}

func isPunct(p rune) bool {
	list := []rune{'+', '-', '*', '/', ';', '<', '>', '(', ')', '='}
	for _, i := range list {
		if i == p {
			return true
		}
	}
	return false
}

func readPunct(text []rune, index int) (int, string) {
	list := [][]rune{
		{'=', '='},
		{'!', '='},
		{'<', '='},
		{'>', '='},
	}
	for _, twoPunct := range list {
		if index+1 < len(text) && text[index] == twoPunct[0] && text[index+1] == twoPunct[1] {
			return 2, string(text[index : index+2])
		}
	}

	if isPunct(text[index]) {
		return 1, string(text[index])
	}

	return 0, ""
}

func tokenize(text []rune) *Token {
	user_input = text

	var head Token
	cur := &head

	index := 0
	for index < len(text) {
		if unicode.IsSpace(text[index]) {
			index++
			continue
		}

		if unicode.IsDigit(text[index]) {
			cur.next = newToken(TK_NUM, index)
			cur = cur.next
			val := getNum(text, &index)
			cur.val = val
			cur.str = strconv.Itoa(val)
			continue
		}

		if unicode.IsLower(text[index]) && unicode.IsLetter(text[index]) {
			cur.next = newToken(TK_VAR, index)
			cur = cur.next
			cur.str = string(text[index])
			index++
			continue
		}

		punctLen, str := readPunct(text, index)
		if punctLen != 0 {
			cur.next = newToken(TK_PUNCT, index)
			cur = cur.next
			cur.str = str
			index += punctLen
			continue
		}

		error_at(index, fmt.Errorf("error tokenize"))
	}

	cur.next = newToken(TK_EOF, index)
	return head.next
}

type NodeKind int

const (
	ND_NUM NodeKind = iota
	ND_ADD
	ND_SUB
	ND_MUL
	ND_DIV
	ND_EQ
	ND_NE
	ND_LT
	ND_LE
	ND_NEG
	ND_VAR
	ND_ASSIGN
	ND_EXPR_STMT
)

type Node struct {
	kind NodeKind
	next *Node
	lhs  *Node
	rhs  *Node
	val  int
	name string
}

func newBinary(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func newNum(val int) *Node {
	return &Node{
		kind: ND_NUM,
		val:  val,
	}
}

func newUnary(kind NodeKind, lhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
	}
}

func newVar(str string) *Node {
	return &Node{
		kind: ND_VAR,
		name: str,
	}
}

func parse(tok *Token) *Node {
	var head Node
	cur := &head

	for tok.kind != TK_EOF {
		cur.next = stmt(&tok, tok)
		cur = cur.next
	}

	return head.next
}

func stmt(rest **Token, tok *Token) *Node {
	node := newUnary(ND_EXPR_STMT, expr(&tok, tok))
	if tok.str != ";" {
		error_at(tok.loc, fmt.Errorf("error stmt"))
	}
	*rest = tok.next
	return node
}

func expr(rest **Token, tok *Token) *Node {
	return assign(rest, tok)
}

func assign(rest **Token, tok *Token) *Node {
	node := equality(&tok, tok)
	if tok.str == "=" {
		node = newBinary(ND_ASSIGN, node, assign(&tok, tok.next))
	}
	*rest = tok
	return node
}

func equality(rest **Token, tok *Token) *Node {
	node := relational(&tok, tok)

	for {
		if tok.str == "==" {
			node = newBinary(ND_EQ, node, relational(&tok, tok.next))
			continue
		}
		if tok.str == "!=" {
			node = newBinary(ND_NE, node, relational(&tok, tok.next))
			continue
		}

		*rest = tok
		return node
	}
}

func relational(rest **Token, tok *Token) *Node {
	node := add(&tok, tok)

	for {
		if tok.str == "<" {
			node = newBinary(ND_LT, node, add(&tok, tok.next))
			continue
		}
		if tok.str == "<=" {
			node = newBinary(ND_LE, node, add(&tok, tok.next))
			continue
		}
		if tok.str == ">" {
			node = newBinary(ND_LT, add(&tok, tok.next), node)
			continue
		}
		if tok.str == ">=" {
			node = newBinary(ND_LE, add(&tok, tok.next), node)
			continue
		}

		*rest = tok
		return node
	}
}

func add(rest **Token, tok *Token) *Node {
	node := mul(&tok, tok)

	for {
		if tok.str == "+" {
			node = newBinary(ND_ADD, node, mul(&tok, tok.next))
			continue
		}

		if tok.str == "-" {
			node = newBinary(ND_SUB, node, mul(&tok, tok.next))
			continue
		}

		*rest = tok
		return node
	}
}

func mul(rest **Token, tok *Token) *Node {
	node := unary(&tok, tok)

	for {
		if tok.str == "*" {
			node = newBinary(ND_MUL, node, unary(&tok, tok.next))
			continue
		}
		if tok.str == "/" {
			node = newBinary(ND_DIV, node, unary(&tok, tok.next))
			continue
		}

		*rest = tok
		return node
	}
}

func unary(rest **Token, tok *Token) *Node {
	if tok.str == "+" {
		return unary(rest, tok.next)
	}

	if tok.str == "-" {
		return newUnary(ND_NEG, unary(rest, tok.next))
	}

	return primary(rest, tok)
}

func primary(rest **Token, tok *Token) *Node {
	if tok.str == "(" {
		node := expr(&tok, tok.next)
		if tok.str != ")" {
			error_at(tok.loc, fmt.Errorf("error primary"))
		}
		*rest = tok.next
		return node
	}

	if tok.kind == TK_NUM {
		node := newNum(tok.val)
		*rest = tok.next
		return node
	}

	if tok.kind == TK_VAR {
		node := newVar(tok.str)
		*rest = tok.next
		return node
	}

	error_at(tok.loc, fmt.Errorf("error primary"))
	return nil
}

var depth int = 0

func push() {
	fmt.Println("	push %rax")
	depth++
}

func pop(p string) {
	fmt.Printf("	pop %s\n", p)
	depth--
}

func gen_addr(node *Node) {
	if node.kind != ND_VAR {
		fmt.Fprintf(os.Stderr, "error gen_addr")
	}

	offset := (rune(node.name[0]) - 'a' + 1) * -8
	fmt.Printf("	lea %d(%%rbp), %%rax\n", offset)
}

func gen_expr(node *Node) {
	switch node.kind {
	case ND_ASSIGN:
		gen_addr(node.lhs)
		push()
		gen_expr(node.rhs)
		pop("%rdi")
		fmt.Println("	mov %rax, (%rdi)")
		return
	case ND_VAR:
		gen_addr(node)
		fmt.Println("	mov (%rax), %rax")
		return
	case ND_NUM:
		fmt.Printf("	mov $%d, %%rax\n", node.val)
		return
	case ND_NEG:
		gen_expr(node.lhs)
		fmt.Println("	neg %rax")
		return
	}

	gen_expr(node.rhs)
	push()
	gen_expr(node.lhs)
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

	fmt.Fprintf(os.Stderr, "error gen_expr %v", node.kind)
}

func codegen(node *Node) {
	fmt.Println(".globl main")
	fmt.Println("main:")

	fmt.Println("	push %rbp")
	fmt.Println("	mov %rsp, %rbp")
	fmt.Println("	sub $208, %rsp")

	for cur := node; cur != nil; cur = cur.next {
		if cur.kind != ND_EXPR_STMT {
			fmt.Fprintf(os.Stderr, "error nd_expr_stmt")
		}
		gen_expr(cur.lhs)
		if depth != 0 {
			fmt.Fprintf(os.Stderr, "Error codegen")
		}
	}

	fmt.Println("	mov %rbp, %rsp")
	fmt.Println("	pop %rbp")
	fmt.Println("	ret")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s: invalid argument\n", os.Args[0])
		os.Exit(1)
	}

	var text []rune
	for _, char := range os.Args[1] {
		text = append(text, char)
	}

	tok := tokenize(text)
	node := parse(tok)
	codegen(node)
}
