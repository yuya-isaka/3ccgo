package main

import (
	"3ccgo5/codegen"
	"3ccgo5/header"
	"3ccgo5/parser"
	"3ccgo5/tokenize"
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
