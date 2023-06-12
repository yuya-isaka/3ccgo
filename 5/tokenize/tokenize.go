package tokenize

import (
	"3ccgo5/header"
	"unicode"
)

func getNum(text []rune, loc *int) int {
	result := int(text[*loc]) - '0'
	*loc++
	for *loc < len(text) && unicode.IsDigit(text[*loc]) {
		result = result*10 + (int(text[*loc]) - '0')
		*loc++
	}
	return result
}

func ispunct(target rune) int {
	list := []rune{
		'=', ';', '+', '-', '*', '/', '(', ')', '<', '>',
	}

	for _, i := range list {
		if target == i {
			return 1
		}
	}

	return 0
}

func readPunct(text []rune, loc int) int {
	if loc+1 < len(text) {
		list := [][]rune{
			{'=', '='},
			{'!', '='},
			{'<', '='},
			{'>', '='},
		}

		for _, chars := range list {
			if text[loc] == chars[0] && text[loc+1] == chars[1] {
				return 2
			}
		}
	}

	return ispunct(text[loc])
}

func newToken(kind header.TokenKind, loc int) *header.Token {
	return &header.Token{
		Kind: kind,
		Loc:  loc,
	}
}

func Tokenize(text []rune) *header.Token {
	header.UserInput = text

	var head header.Token
	cur := &head

	loc := 0
	for loc < len(text) {
		if unicode.IsSpace(text[loc]) {
			loc++
			continue
		}

		if unicode.IsDigit(text[loc]) {
			cur.Next = newToken(header.TK_NUM, loc)
			cur = cur.Next
			cur.Val = getNum(text, &loc)
			continue
		}

		punctLen := readPunct(text, loc)
		if punctLen != 0 {
			cur.Next = newToken(header.TK_PUNCT, loc)
			cur = cur.Next
			cur.Name = string(text[loc : loc+punctLen])
			loc += punctLen
			continue
		}

		if unicode.IsLower(text[loc]) && unicode.IsLetter(text[loc]) {
			cur.Next = newToken(header.TK_VAR, loc)
			cur = cur.Next
			cur.Name = string(text[loc])
			loc++
			continue
		}

		header.ErrorAt(loc, "error tokenize")
	}

	cur.Next = newToken(header.TK_EOF, loc-1)
	return head.Next
}
