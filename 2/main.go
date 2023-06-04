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

func tokenize(text []rune) *Token {
	user_input = text

	head := Token{}
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

		if text[index] == '+' || text[index] == '-' {
			cur.next = newToken(TK_PUNCT, index)
			cur = cur.next
			cur.str = string(text[index])
			index++
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
)

type Node struct {
	kind NodeKind
	next *Node
	lhs  *Node
	rhs  *Node
	val  int
}

// func expr(tok *Token) (*Node, *Token) {

// }

// func parse(tok *Token) *Node {
// 	var node *Node
// 	node, tok = expr(tok)
// }

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s: invalid argument", os.Args[0])
	}

	text := []rune{}
	for _, char := range os.Args[1] {
		text = append(text, char)
	}

	tok := tokenize(text)
	// node := parse(tok)

	fmt.Println(".globl main")
	fmt.Println("main:")
	fmt.Printf("	mov $%d, %%rax\n", tok.val)
	tok = tok.next

	for tok.kind != TK_EOF {
		if tok.kind == TK_NUM {
			fmt.Printf("	mov $%d, %%rax\n", tok.val)
			tok = tok.next
			continue
		} else if tok.kind == TK_PUNCT {
			if tok.str == "+" {
				fmt.Printf("	add $%d, %%rax\n", tok.next.val)
				tok = tok.next.next
				continue
			} else if tok.str == "-" {
				fmt.Printf("	sub $%d, %%rax\n", tok.next.val)
				tok = tok.next.next
				continue
			}
		} else {
			error_at(tok.loc, fmt.Errorf("error main"))
		}
	}

	fmt.Println("	ret")
}
