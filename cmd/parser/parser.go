package parser

import (
	"fmt"
	"strings"

	. "github.com/nathanleiby/glox/cmd/models"
)

type Parser struct {
	tokens  []Token
	currIdx int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
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

func (expr *BinaryExpr) String() string {
	return fmt.Sprintf("%s %s %s", expr.Left, expr.Operator.TokenType, expr.Right)
}

func (expr *BinaryExpr) Expr() {
	return
}

// UnaryExpr is an unary expression
type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (expr *UnaryExpr) String() string {
	return fmt.Sprintf("%s %s", expr.Operator, expr.Right)
}

func (expr *UnaryExpr) Expr() {
	return
}

// LiteralExpr is an unary expression
type LiteralExpr struct {
	Value interface{}
}

func (expr *LiteralExpr) String() string {
	return fmt.Sprintf("%v", expr.Value)
}

func (expr *LiteralExpr) Expr() {
	return
}

// GroupingExpr is an unary expression
type GroupingExpr struct {
	Value Expr
}

func (expr *GroupingExpr) String() string {
	return fmt.Sprintf("(%v)", expr.Value)
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
	curType := p.current().TokenType
	return curType == tt
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

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	var expr Expr
	var err error
	expr, err = p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Operator: op,
			Right:    right,
		}, nil
	}

	return p.primary()
}

var ErrExpectExpression = fmt.Errorf("expect expression")

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &LiteralExpr{Value: false}, nil
	}
	if p.match(TRUE) {
		return &LiteralExpr{Value: true}, nil
	}
	if p.match(NIL) {
		return &LiteralExpr{Value: nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{Value: p.previous().Literal}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &GroupingExpr{Value: expr}, nil
	}

	return nil, ErrExpectExpression
}

func (p *Parser) consume(tt TokenType, message string) (Token, error) {
	if p.check(tt) {
		return p.advance(), nil
	}

	return Token{}, parseError(p.current(), message)
}

var ErrParsingEOF = fmt.Errorf("error parsing, at EOF")     // TODO: pass context
var ErrParsingUnknown = fmt.Errorf("error parsing, unkown") // TODO: pass context

func parseError(t Token, msg string) error {
	if t.TokenType == EOF {
		// report(t.Line, " at end", msg)
		return ErrParsingEOF
	} else {
		// report(t.Line, " at '"+t.Lexeme+"'", msg)
		return ErrParsingUnknown
	}
}

// TODO: refactor to an erros package
func report(line int, where, message string) {
	// FUTURE: Add more useful error context
	//
	// Error: Unexpected "," in argument list.

	// 15 | function(first, second,);
	//                            ^-- Here.
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

// synchronizes the token stream
// discards tokens until it reaches one that can appear at that point in the rule.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}

		switch p.current().TokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}

func Parenthesize(exprs []Expr) string {
	s := make([]string, len(exprs))
	for i, expr := range exprs {
		s[i] = fmt.Sprintf("%s", expr)
	}

	return fmt.Sprintf("(%s)", strings.Join(s, " "))
}
