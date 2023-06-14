package main

import (
	"3ccgo6/codegen"
	"3ccgo6/header"
	"3ccgo6/parser"
	"3ccgo6/tokenize"
	"fmt"
	"os"
)

var index int = 0

func lhsrhsDebug(node *header.Node) {
	if node.Lhs != nil {
		fmt.Printf("%*s", index, "")

		fmt.Printf("Lhs: %v", node.Lhs.Kind)
		if node.Lhs.Kind == header.NdNum {
			fmt.Printf(" → %v", node.Lhs.Val)
		} else if node.Lhs.Kind == header.NdVar {
			fmt.Printf(" → %v", node.Lhs.Var.Name)
		}
		fmt.Println("")

		index += 2
		lhsrhsDebug(node.Lhs)
		index -= 2
	}

	if node.Rhs != nil {
		fmt.Printf("%*s", index, "")

		fmt.Printf("Rhs: %v", node.Rhs.Kind)
		if node.Rhs.Kind == header.NdNum {
			fmt.Printf(" → %v", node.Rhs.Val)
		} else if node.Rhs.Kind == header.NdVar {
			fmt.Printf(" → %v", node.Rhs.Var.Name)
		}
		fmt.Println("")

		index += 2
		lhsrhsDebug(node.Rhs)
		index -= 2
	}
}

func main() {
	if len(os.Args) != 2 {
		header.Errorf("%s invalid argument", os.Args[0])
	}

	var text []rune
	for _, char := range os.Args[1] {
		text = append(text, char)
	}

	tok := tokenize.Tokenize(text)

	// // Debug ----------------------------------------------------------------------
	// for cur := tok; cur != nil; cur = cur.Next {
	// 	fmt.Printf("kind: %v \t val: %v \t Name: %v\n", cur.Kind, cur.Val, cur.Name)
	// }
	// // ----------------------------------------------------------------------------

	function := parser.Parser(tok)

	// // Debug ----------------------------------------------------------------------
	// for cur := function.Body; cur != nil; cur = cur.Next {
	// 	fmt.Printf("\n ---Kind: %v---\n", cur.Kind)

	// 	lhsrhsDebug(cur)
	// }
	// fmt.Println("")
	// // ----------------------------------------------------------------------------

	codegen.Codegen(function)
}
