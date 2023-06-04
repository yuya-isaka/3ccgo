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
)

type Token struct {
	kind TokenKind
	loc  int
	val  int
	str  string
	next *Token
}

var user_input string

func errorDefault(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func errorAt(loc int, err error) {
	fmt.Fprintf(os.Stderr, "%s\n", user_input)
	fmt.Fprintf(os.Stderr, "%*s", loc, "")
	fmt.Fprintf(os.Stderr, "^ ")
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func newToken(kind TokenKind, loc int, val int, str string) *Token {
	tok := &Token{kind: kind, loc: loc, val: val, str: str}
	return tok
}

func tokenize(input string) *Token {
	user_input = input

	head := &Token{}
	cur := head

	for i := 0; i < len(input); i++ {
		if unicode.IsSpace(rune(input[i])) {
			continue
		}

		if input[i] == '+' || input[i] == '-' {
			cur.next = newToken(TK_PUNCT, i, 0, string(input[i]))
			cur = cur.next
			continue
		}

		if unicode.IsDigit(rune(input[i])) {
			j := i
			for j < len(input) && unicode.IsDigit(rune(input[j])) {
				j++
			}
			val, err := strconv.Atoi(input[i:j])
			if err != nil {
				errorAt(i, err)
			}
			cur.next = newToken(TK_NUM, i, val, input[i:j])
			cur = cur.next
			i = j - 1
			continue
		}

		errorAt(i, fmt.Errorf("invalid token"))
	}

	cur.next = &Token{kind: TK_EOF, loc: len(input)}
	return head.next
}

func main() {
	if len(os.Args) != 2 {
		errorDefault(fmt.Errorf("%s: incorrect number of arguments.", os.Args[0]))
	}

	input := os.Args[1]
	tok := tokenize(input)

	fmt.Println(".globl main")
	fmt.Println("main:")
	fmt.Printf("	mov $%d, %%rax\n", tok.val)
	tok = tok.next

	for tok.kind != TK_EOF {
		if tok.str == "+" {
			tok = tok.next
			fmt.Printf("	add $%d, %%rax\n", tok.val)
		} else if tok.str == "-" {
			tok = tok.next
			fmt.Printf("	sub $%d, %%rax\n", tok.val)
		} else {
			errorAt(tok.loc, fmt.Errorf("expected '+' or '-' but got %v", tok.val))
		}
		tok = tok.next
	}

	fmt.Println("	ret")
}
