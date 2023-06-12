package parser

import "3ccgo5/header"

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

func stmt(tok *header.Token) (*header.Node, *header.Token) {
	var lhs *header.Node
	lhs, tok = expr(tok)
	node := newUnary(header.ND_EXPR_STMT, lhs)
	if tok.Name != ";" {
		header.ErrorAt(tok.Loc, "error stmt")
	}
	return node, tok.Next
}

func expr(tok *header.Token) (*header.Node, *header.Token) {
	return assign(tok)
}

func assign(tok *header.Token) (*header.Node, *header.Token) {
	var node *header.Node
	node, tok = equivalent(tok)

	if tok.Name == "=" {
		var rhs *header.Node
		rhs, tok = assign(tok.Next)
		node = newBinary(header.ND_ASSIGN, node, rhs)
	}
	return node, tok
}

func equivalent(tok *header.Token) (*header.Node, *header.Token) {
	var node *header.Node
	node, tok = relation(tok)

	for {
		var rhs *header.Node

		if tok.Name == "==" {
			rhs, tok = relation(tok.Next)
			node = newBinary(header.ND_EQ, node, rhs)
			continue
		}

		if tok.Name == "!=" {
			rhs, tok = relation(tok.Next)
			node = newBinary(header.ND_NE, node, rhs)
			continue
		}

		return node, tok
	}
}

func relation(tok *header.Token) (*header.Node, *header.Token) {
	var node *header.Node
	node, tok = add(tok)

	for {
		var rhs *header.Node

		if tok.Name == "<" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.ND_LT, node, rhs)
			continue
		}

		if tok.Name == "<=" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.ND_LE, node, rhs)
			continue
		}

		if tok.Name == ">" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.ND_LT, rhs, node)
			continue
		}

		if tok.Name == ">=" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.ND_LE, rhs, node)
			continue
		}

		return node, tok
	}
}

func add(tok *header.Token) (*header.Node, *header.Token) {
	var node *header.Node
	node, tok = mul(tok)

	for {
		var rhs *header.Node

		if tok.Name == "+" {
			rhs, tok = mul(tok.Next)
			node = newBinary(header.ND_ADD, node, rhs)
			continue
		}

		if tok.Name == "-" {
			rhs, tok = mul(tok.Next)
			node = newBinary(header.ND_SUB, node, rhs)
			continue
		}

		return node, tok
	}
}

func mul(tok *header.Token) (*header.Node, *header.Token) {
	var node *header.Node
	node, tok = unary(tok)

	for {
		var rhs *header.Node

		if tok.Name == "*" {
			rhs, tok = unary(tok.Next)
			node = newBinary(header.ND_MUL, node, rhs)
			continue
		}

		if tok.Name == "/" {
			rhs, tok = unary(tok.Next)
			node = newBinary(header.ND_DIV, node, rhs)
			continue
		}

		return node, tok
	}
}

func unary(tok *header.Token) (*header.Node, *header.Token) {
	if tok.Name == "+" {
		return unary(tok.Next)
	}

	if tok.Name == "-" {
		var node *header.Node
		node, tok = unary(tok.Next)

		return newUnary(header.ND_NEG, node), tok
	}

	return primary(tok)
}

func primary(tok *header.Token) (*header.Node, *header.Token) {
	if tok.Name == "(" {
		var node *header.Node
		node, tok = expr(tok.Next)
		if tok.Name != ")" {
			header.ErrorAt(tok.Loc, "error primary")
		}

		return node, tok.Next
	}

	if tok.Kind == header.TK_NUM {
		return newNum(tok.Val), tok.Next
	}

	if tok.Kind == header.TK_VAR {
		return newVar(tok.Name), tok.Next
	}

	header.ErrorAt(tok.Loc, "error primary end")
	return nil, nil
}

func Parse(tok *header.Token) *header.Node {
	var head header.Node
	cur := &head

	for tok.Kind != header.TK_EOF {
		cur.Next, tok = stmt(tok)
		cur = cur.Next
	}

	return head.Next
}
