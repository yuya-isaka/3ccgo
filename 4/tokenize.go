package main

import (
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

func Tokenize(text []rune) *Token {
	UserInput = text

	var head Token
	cur := &head

	loc := 0
	for loc < len(text) {
		if unicode.IsSpace(text[loc]) {
			loc++
			continue
		}

		if unicode.IsDigit(text[loc]) {
			cur.next = newToken(TK_NUM, loc)
			cur = cur.next
			cur.val = getNum(text, &loc)
			continue
		}

		punctLen := readPunct(text, loc)
		if punctLen != 0 {
			cur.next = newToken(TK_PUNCT, loc)
			cur = cur.next
			cur.name = string(text[loc : loc+punctLen])
			loc += punctLen
			continue
		}

		if unicode.IsLower(text[loc]) && unicode.IsLetter(text[loc]) {
			cur.next = newToken(TK_VAR, loc)
			cur = cur.next
			cur.name = string(text[loc])
			loc++
			continue
		}

		ErrorAt(loc, "error tokenize")
	}

	cur.next = newToken(TK_EOF, loc-1)
	return head.next
}
