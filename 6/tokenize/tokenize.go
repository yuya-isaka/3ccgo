package tokenize

import (
	"3ccgo6/header"
	"unicode"
)

func newToken(kind header.TokenKind, i int) *header.Token {
	return &header.Token{
		Kind: kind,
		Loc:  i,
	}
}

func getNum(text []rune, i int) (int, int) {
	result := int(text[i]) - '0'
	i++
	for i < len(text) && unicode.IsDigit(text[i]) {
		result = 10*result + (int(text[i]) - '0')
		i++
	}
	return result, i
}

func isPunct(target rune) int {
	list := []rune{
		'+', '-', '*', '/', '=', ';', '<', '>', '(', ')',
	}
	for _, li := range list {
		if target == li {
			return 1
		}
	}

	return 0
}

func readPunct(text []rune, i int) int {
	if i+1 < len(text) {
		list := [][]rune{
			{'=', '='},
			{'!', '='},
			{'<', '='},
			{'>', '='},
		}
		for _, li := range list {
			if text[i] == li[0] && text[i+1] == li[1] {
				return 2
			}
		}
	}

	return isPunct(text[i])
}

func Tokenize(text []rune) *header.Token {
	header.Userinput = text

	var head header.Token
	cur := &head

	i := 0
	for i < len(text) {
		if unicode.IsSpace(text[i]) {
			i++
			continue
		}

		if unicode.IsDigit(text[i]) {
			cur.Next = newToken(header.TkNum, i)
			cur = cur.Next
			cur.Val, i = getNum(text, i)
			continue
		}

		punctLen := readPunct(text, i)
		if punctLen != 0 {
			cur.Next = newToken(header.TkPunct, i)
			cur = cur.Next
			cur.Name = string(text[i : i+punctLen])
			i += punctLen
			continue
		}

		if unicode.IsLower(text[i]) && unicode.IsLetter(text[i]) {
			cur.Next = newToken(header.TkVar, i)
			cur = cur.Next
			cur.Name = string(text[i])
			i++
			continue
		}

		header.ErrorAt(i, "error tokenize")
	}

	cur.Next = newToken(header.TkEof, i)
	return head.Next
}
