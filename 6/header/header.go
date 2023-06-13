package header

import (
	"fmt"
	"os"
)

var Userinput []rune

func Errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf(format, a...))
	os.Exit(1)
}

func ErrorAt(loc int, msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", string(Userinput))
	fmt.Fprintf(os.Stderr, "%*s", loc, "")
	fmt.Fprintf(os.Stderr, "^ ")
	Errorf(msg)
}

type TokenKind int

const (
	TkPunct TokenKind = iota
	TkVar
	TkNum
	TkEof
)

func (tk TokenKind) String() string {
	switch tk {
	case TkPunct:
		return "Punct"
	case TkVar:
		return "Var"
	case TkNum:
		return "Num"
	case TkEof:
		return "EOF"
	}

	return ""
}

type Token struct {
	Kind TokenKind
	Next *Token
	Loc  int
	Val  int
	Name string
}

type NodeKind int

const (
	NdAdd NodeKind = iota
	NdSub
	NdMul
	NdDiv
	NdEq
	NdNe
	NdLt
	NdLe
	NdVar
	NdAssign
	NdExprStmt
	NdNum
	NdNeg
)

func (nk NodeKind) String() string {
	switch nk {
	case NdAdd:
		return "Add"
	case NdSub:
		return "Sub"
	case NdMul:
		return "Mul"
	case NdDiv:
		return "Div"
	case NdEq:
		return "Eq"
	case NdNe:
		return "Ne"
	case NdLt:
		return "Lt"
	case NdLe:
		return "Le"
	case NdAssign:
		return "Assign"
	case NdExprStmt:
		return "Expr Stmt"
	case NdVar:
		return "Var"
	case NdNum:
		return "Num"
	case NdNeg:
		return "Neg"
	}

	return ""
}

type Node struct {
	Kind NodeKind
	Next *Node
	Lhs  *Node
	Rhs  *Node
	Val  int
	Name string
}
