package cmd

type Parser struct {
	tokens  []Token
	currIdx int
}

// Expr is an expression
type Expr interface {
	Expr() // TODO
}

// BinaryExpr is a binary expression
type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr *BinaryExpr) Expr() {
	return
}

// UnaryExpr is an unary expression
type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (expr *UnaryExpr) Expr() {
	return
}

////////////////////////////////////////
// Inspecting position
////////////////////////////////////////
// current is called peek() in the book
func (p *Parser) current() Token {
	return p.tokens[p.currIdx]
}

func (p *Parser) previous() Token {
	return p.tokens[p.currIdx-1]
}

func (p *Parser) isAtEnd() bool {
	return p.current().tokenType == EOF
}

func (p *Parser) check(tt TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.current().tokenType == tt
}

////////////////////////////////////////
// Movement
////////////////////////////////////////
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.currIdx++
	}
	return p.previous()
}

// match is "advance if check"
func (p *Parser) match(tts ...TokenType) bool {
	for _, tt := range tts {
		if p.check(tt) {
			p.advance()
			return true
		}
	}
	return false
}

////////////////////////////////////////
// Expressions
////////////////////////////////////////
func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		op := p.previous()
		right := p.factor()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	var expr Expr
	expr = p.unary()

	for p.match([]TokenType{SLASH, STAR}) {
		op := p.previous()
		right := p.unary()
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match([]TokenType{BANG, MINUS}) {
		op := p.previous()
		right := p.unary()
		return &UnaryExpr{
			Operator: op,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &LiteralExpr{value: false}
	}
	if p.match(TRUE) {
		return &LiteralExpr{value: true}
	}
	if p.match(NIL) {
		return &LiteralExpr{value: nil}
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{value: p.previous().literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{group: expr}
	}
}
