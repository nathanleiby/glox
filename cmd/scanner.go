package cmd

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   string // TODO: Object?
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", string(t.tokenType), t.lexeme, t.literal)
}

type Scanner struct {
	source  string
	start   int
	current int
	line    int

	tokens []Token
}

func (s Scanner) Tokens() []Token {
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s Scanner) advance() byte {
	s.current++
	return s.source[s.current]
}

func (s Scanner) addToken(tt TokenType) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		tokenType: tt,
		lexeme:    text,
		line:      s.line,
	})

}

func (s Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	default:
		loxError(s.line, fmt.Sprintf("unexpected character: %c", c))
	}

}

func (s Scanner) ScanTokens() error {
	s.start = 0
	s.current = 0
	s.line = 1

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		tokenType: EOF,
		lexeme:    "",
		line:      s.line,
	})

	return nil
}

func NewScanner(source string) Scanner {
	return Scanner{
		source: source,
	}
}
