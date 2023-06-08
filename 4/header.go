package main

import (
	"fmt"
	"os"
)

var UserInput []rune

func Errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", fmt.Errorf(format, a...))
	os.Exit(1)
}

func ErrorAt(loc int, msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", string(UserInput))
	fmt.Fprintf(os.Stderr, "%*s", loc, "")
	fmt.Fprintf(os.Stderr, "^ ")
	Errorf(msg)
}

func newToken(kind TokenKind, loc int) *Token {
	return &Token{
		kind: kind,
		loc:  loc,
	}
}

type TokenKind int

const (
	TK_PUNCT TokenKind = iota
	TK_VAR
	TK_NUM
	TK_EOF
)

type Token struct {
	kind TokenKind
	next *Token
	loc  int
	val  int
	name string
}

type NodeKind int

const (
	ND_ADD NodeKind = iota
	ND_SUB
	ND_MUL
	ND_DIV
	ND_EQ
	ND_NE
	ND_LT
	ND_LE
	ND_NUM
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
