package main

import (
	"3ccgo4/codegen"
	"3ccgo4/header"
	"3ccgo4/parser"
	"3ccgo4/tokenize"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		header.Errorf("%s: invalid argument", os.Args[0])
	}

	var text []rune
	for _, i := range os.Args[1] {
		text = append(text, i)
	}

	tok := tokenize.Tokenize(text)
	node := parser.Parse(tok)
	codegen.Codegen(node)
}
