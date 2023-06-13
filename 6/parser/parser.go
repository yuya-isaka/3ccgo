package parser

import "3ccgo6/header"

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
		Kind: header.NdNum,
		Val:  val,
	}
}

func newVar(name string) *header.Node {
	return &header.Node{
		Kind: header.NdVar,
		Name: name,
	}
}

func stmt(tok *header.Token) (*header.Node, *header.Token) {
	node, tok := expr(tok)
	node = newUnary(header.NdExprStmt, node)
	if tok.Name != ";" {
		header.ErrorAt(tok.Loc, "Error stmt")
	}

	return node, tok.Next
}

func expr(tok *header.Token) (*header.Node, *header.Token) {
	return assign(tok)
}

func assign(tok *header.Token) (*header.Node, *header.Token) {
	node, tok := equivalent(tok)

	if tok.Name == "=" {
		var rhs *header.Node
		rhs, tok = assign(tok.Next)
		node = newBinary(header.NdAssign, node, rhs)
	}

	return node, tok
}

func equivalent(tok *header.Token) (*header.Node, *header.Token) {
	node, tok := relation(tok)

	for {
		var rhs *header.Node

		if tok.Name == "==" {
			rhs, tok = relation(tok.Next)
			node = newBinary(header.NdEq, node, rhs)
			continue
		}

		if tok.Name == "!=" {
			rhs, tok = relation(tok.Next)
			node = newBinary(header.NdNe, node, rhs)
			continue
		}

		return node, tok
	}
}

func relation(tok *header.Token) (*header.Node, *header.Token) {
	node, tok := add(tok)

	for {
		var rhs *header.Node

		if tok.Name == "<" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.NdLt, node, rhs)
			continue
		}

		if tok.Name == "<=" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.NdLe, node, rhs)
			continue
		}

		if tok.Name == ">" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.NdLt, rhs, node)
			continue
		}

		if tok.Name == ">=" {
			rhs, tok = add(tok.Next)
			node = newBinary(header.NdLe, rhs, node)
			continue
		}

		return node, tok
	}
}

func add(tok *header.Token) (*header.Node, *header.Token) {
	node, tok := mul(tok)

	for {
		var rhs *header.Node

		if tok.Name == "+" {
			rhs, tok = mul(tok.Next)
			node = newBinary(header.NdAdd, node, rhs)
			continue
		}

		if tok.Name == "-" {
			rhs, tok = mul(tok.Next)
			node = newBinary(header.NdSub, node, rhs)
			continue
		}

		return node, tok
	}
}

func mul(tok *header.Token) (*header.Node, *header.Token) {
	node, tok := unary(tok)

	for {
		var rhs *header.Node

		if tok.Name == "*" {
			rhs, tok = unary(tok.Next)
			node = newBinary(header.NdMul, node, rhs)
			continue
		}

		if tok.Name == "/" {
			rhs, tok = unary(tok.Next)
			node = newBinary(header.NdDiv, node, rhs)
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
		node, tok := unary(tok.Next)
		node = newUnary(header.NdNeg, node)
		return node, tok
	}

	return primary(tok)
}

func primary(tok *header.Token) (*header.Node, *header.Token) {
	if tok.Name == "(" {
		node, tok := expr(tok.Next)
		if tok.Name != ")" {
			header.ErrorAt(tok.Loc, "Error primary")
		}
		return node, tok.Next
	}

	if tok.Kind == header.TkNum {
		return newNum(tok.Val), tok.Next
	}

	if tok.Kind == header.TkVar {
		return newVar(tok.Name), tok.Next
	}

	header.ErrorAt(tok.Loc, "Error primary last")
	return nil, nil
}

func Parser(tok *header.Token) *header.Node {
	var head header.Node
	cur := &head

	for tok.Kind != header.TkEof {
		cur.Next, tok = stmt(tok)
		cur = cur.Next
	}

	return head.Next
}
