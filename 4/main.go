package main

import (
	"os"
)

func main() {
	if len(os.Args) != 2 {
		Errorf("%s: invalid argument", os.Args[0])
	}

	var text []rune
	for _, i := range os.Args[1] {
		text = append(text, i)
	}

	tok := Tokenize(text)
	node := Parse(tok)
	Codegen(node)
}
