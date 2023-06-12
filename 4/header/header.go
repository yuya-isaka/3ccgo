package header

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

type TokenKind int

const (
	TK_PUNCT TokenKind = iota
	TK_VAR
	TK_NUM
	TK_EOF
)

type Token struct {
	Kind TokenKind
	Next *Token
	Loc  int
	Val  int
	Name string
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
	Kind NodeKind
	Next *Node
	Lhs  *Node
	Rhs  *Node
	Val  int
	Name string
}
