package parser

import (
	. "github.com/nathanleiby/glox/cmd/models"
)

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

// LiteralExpr is an unary expression
type LiteralExpr struct {
	Value interface{}
}

func (expr *LiteralExpr) Expr() {
	return
}

// GroupingExpr is an unary expression
type GroupingExpr struct {
	Value Expr
}

func (expr *GroupingExpr) Expr() {
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
	return p.current().TokenType == EOF
}

func (p *Parser) check(tt TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.current().TokenType == tt
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
		return &LiteralExpr{Value: false}
	}
	if p.match(TRUE) {
		return &LiteralExpr{Value: true}
	}
	if p.match(NIL) {
		return &LiteralExpr{Value: nil}
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{Value: p.previous().Literal}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &GroupingExpr{Value: expr}
	}

	panic("primary() failed to match anything")
}

func (p *Parser) consume(tt TokenType, message string) Token {
	if p.check(tt) {
		return p.advance()
	}

	panic(p.parseError(p.current(), message))
}

func (p *Parser) parseError(t Token, msg string) {
	if t.TokenType == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}
