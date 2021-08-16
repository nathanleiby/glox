package parser

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
