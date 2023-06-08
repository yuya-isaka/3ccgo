package main

func newUnary(kind NodeKind, lhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
	}
}

func newBinary(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func newNum(val int) *Node {
	return &Node{
		kind: ND_NUM,
		val:  val,
	}
}

func newVar(name string) *Node {
	return &Node{
		kind: ND_VAR,
		name: name,
	}
}

func stmt(rest **Token, tok *Token) *Node {
	node := newUnary(ND_EXPR_STMT, expr(&tok, tok))
	if tok.name != ";" {
		ErrorAt(tok.loc, "error stmt")
	}
	*rest = tok.next
	return node
}

func expr(rest **Token, tok *Token) *Node {
	return assign(rest, tok)
}

func assign(rest **Token, tok *Token) *Node {
	node := equivalent(&tok, tok)
	if tok.name == "=" {
		node = newBinary(ND_ASSIGN, node, assign(&tok, tok.next))
	}
	*rest = tok
	return node
}

func equivalent(rest **Token, tok *Token) *Node {
	node := relation(&tok, tok)

	for {
		if tok.name == "==" {
			node = newBinary(ND_EQ, node, relation(&tok, tok.next))
			continue
		}

		if tok.name == "!=" {
			node = newBinary(ND_NE, node, relation(&tok, tok.next))
			continue
		}

		*rest = tok
		return node
	}
}

func relation(rest **Token, tok *Token) *Node {
	node := add(&tok, tok)

	for {
		if tok.name == "<" {
			node = newBinary(ND_LT, node, add(&tok, tok.next))
			continue
		}

		if tok.name == "<=" {
			node = newBinary(ND_LE, node, add(&tok, tok.next))
			continue
		}

		if tok.name == ">" {
			node = newBinary(ND_LT, add(&tok, tok.next), node)
			continue
		}

		if tok.name == ">=" {
			node = newBinary(ND_LE, add(&tok, tok.next), node)
			continue
		}

		*rest = tok
		return node
	}
}

func add(rest **Token, tok *Token) *Node {
	node := mul(&tok, tok)

	for {
		if tok.name == "+" {
			node = newBinary(ND_ADD, node, mul(&tok, tok.next))
			continue
		}

		if tok.name == "-" {
			node = newBinary(ND_SUB, node, mul(&tok, tok.next))
			continue
		}

		*rest = tok
		return node
	}
}

func mul(rest **Token, tok *Token) *Node {
	node := unary(&tok, tok)

	for {
		if tok.name == "*" {
			node = newBinary(ND_MUL, node, unary(&tok, tok.next))
			continue
		}

		if tok.name == "/" {
			node = newBinary(ND_DIV, node, unary(&tok, tok.next))
			continue
		}

		*rest = tok
		return node
	}
}

func unary(rest **Token, tok *Token) *Node {
	if tok.name == "+" {
		return unary(rest, tok.next)
	}

	if tok.name == "-" {
		return newUnary(ND_NEG, unary(rest, tok.next))
	}

	return primary(rest, tok)
}

func primary(rest **Token, tok *Token) *Node {
	if tok.name == "(" {
		node := expr(&tok, tok.next)
		if tok.name != ")" {
			ErrorAt(tok.loc, "error primary")
		}
		*rest = tok.next
		return node
	}

	if tok.kind == TK_NUM {
		node := newNum(tok.val)
		*rest = tok.next
		return node
	}

	if tok.kind == TK_VAR {
		node := newVar(tok.name)
		*rest = tok.next
		return node
	}

	ErrorAt(tok.loc, "error primary end")
	return nil
}

func Parse(tok *Token) *Node {
	var head Node
	cur := &head

	for tok.kind != TK_EOF {
		cur.next = stmt(&tok, tok)
		cur = cur.next
	}

	return head.next
}
