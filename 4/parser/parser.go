package parser

import "3ccgo4/header"

func newUnary(kind header.NodeKind, lhs *header.Node) *header.Node {
	return &header.Node{
		Kind: kind,
		Lhs:  lhs,
	}
}

func newBinary(kind header.NodeKind, lhs *header.Node, rhs *header.Node) *header.Node {
	return &header.Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func newNum(val int) *header.Node {
	return &header.Node{
		Kind: header.ND_NUM,
		Val:  val,
	}
}

func newVar(name string) *header.Node {
	return &header.Node{
		Kind: header.ND_VAR,
		Name: name,
	}
}

func stmt(rest **header.Token, tok *header.Token) *header.Node {
	node := newUnary(header.ND_EXPR_STMT, expr(&tok, tok))
	if tok.Name != ";" {
		header.ErrorAt(tok.Loc, "error stmt")
	}
	*rest = tok.Next
	return node
}

func expr(rest **header.Token, tok *header.Token) *header.Node {
	return assign(rest, tok)
}

func assign(rest **header.Token, tok *header.Token) *header.Node {
	node := equivalent(&tok, tok)
	if tok.Name == "=" {
		node = newBinary(header.ND_ASSIGN, node, assign(&tok, tok.Next))
	}
	*rest = tok
	return node
}

func equivalent(rest **header.Token, tok *header.Token) *header.Node {
	node := relation(&tok, tok)

	for {
		if tok.Name == "==" {
			node = newBinary(header.ND_EQ, node, relation(&tok, tok.Next))
			continue
		}

		if tok.Name == "!=" {
			node = newBinary(header.ND_NE, node, relation(&tok, tok.Next))
			continue
		}

		*rest = tok
		return node
	}
}

func relation(rest **header.Token, tok *header.Token) *header.Node {
	node := add(&tok, tok)

	for {
		if tok.Name == "<" {
			node = newBinary(header.ND_LT, node, add(&tok, tok.Next))
			continue
		}

		if tok.Name == "<=" {
			node = newBinary(header.ND_LE, node, add(&tok, tok.Next))
			continue
		}

		if tok.Name == ">" {
			node = newBinary(header.ND_LT, add(&tok, tok.Next), node)
			continue
		}

		if tok.Name == ">=" {
			node = newBinary(header.ND_LE, add(&tok, tok.Next), node)
			continue
		}

		*rest = tok
		return node
	}
}

func add(rest **header.Token, tok *header.Token) *header.Node {
	node := mul(&tok, tok)

	for {
		if tok.Name == "+" {
			node = newBinary(header.ND_ADD, node, mul(&tok, tok.Next))
			continue
		}

		if tok.Name == "-" {
			node = newBinary(header.ND_SUB, node, mul(&tok, tok.Next))
			continue
		}

		*rest = tok
		return node
	}
}

func mul(rest **header.Token, tok *header.Token) *header.Node {
	node := unary(&tok, tok)

	for {
		if tok.Name == "*" {
			node = newBinary(header.ND_MUL, node, unary(&tok, tok.Next))
			continue
		}

		if tok.Name == "/" {
			node = newBinary(header.ND_DIV, node, unary(&tok, tok.Next))
			continue
		}

		*rest = tok
		return node
	}
}

func unary(rest **header.Token, tok *header.Token) *header.Node {
	if tok.Name == "+" {
		return unary(rest, tok.Next)
	}

	if tok.Name == "-" {
		return newUnary(header.ND_NEG, unary(rest, tok.Next))
	}

	return primary(rest, tok)
}

func primary(rest **header.Token, tok *header.Token) *header.Node {
	if tok.Name == "(" {
		node := expr(&tok, tok.Next)
		if tok.Name != ")" {
			header.ErrorAt(tok.Loc, "error primary")
		}
		*rest = tok.Next
		return node
	}

	if tok.Kind == header.TK_NUM {
		node := newNum(tok.Val)
		*rest = tok.Next
		return node
	}

	if tok.Kind == header.TK_VAR {
		node := newVar(tok.Name)
		*rest = tok.Next
		return node
	}

	header.ErrorAt(tok.Loc, "error primary end")
	return nil
}

func Parse(tok *header.Token) *header.Node {
	var head header.Node
	cur := &head

	for tok.Kind != header.TK_EOF {
		cur.Next = stmt(&tok, tok)
		cur = cur.Next
	}

	return head.Next
}
