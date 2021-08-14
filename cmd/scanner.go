package cmd

import "fmt"

type Scanner struct {
	source  string
	start   int
	current int
	line    int

	tokens []Token
}

func (s *Scanner) Tokens() []Token {
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	val := s.source[s.current]
	s.current = s.current + 1
	return val
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// match is a "conditional advance()".
// It checks if the current character matches, and advances if so/
func (s *Scanner) match(expected byte) bool {
	if s.peek() != expected {
		return false
	}

	s.current = s.current + 1
	return true
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		loxError(s.line, "Unterminated string.")
		return
	}

	s.advance() // Consume the closing "

	// Trim the surrounding quotes.
	// text := s.tokenText()
	// value := text[1 : len(text)-2]
	s.addToken(STRING) //  TODO: addTokenliteral with real value
}

func (s *Scanner) scanNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// look for a fractional part
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// consume the '.'
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	// TODO: Support adding a number token
	s.addToken(NUMBER)
}

func (s *Scanner) tokenText() string {
	return s.source[s.start:s.current]
}

func (s *Scanner) scanIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.tokenText()
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) addToken(tt TokenType) {
	text := s.tokenText()
	s.tokens = append(s.tokens, Token{
		tokenType: tt,
		lexeme:    text,
		line:      s.line,
	})
}

// // TODO: Support non-strings, like numbers/etc
// // the original code has an Object in java. Maybe interface{} could work here?
// func (s *Scanner) addTokenLiteral(tt TokenType, value string) {
//     text = source.substring(start, current);
//     tokens.add(new Token(type, text, literal, line));
// }

func (s *Scanner) scanToken() {
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
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	// Whitespace
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line = s.line + 1
	case '"':
		s.scanString()
	default:
		if s.isDigit(c) {
			s.scanNumber()
		} else if s.isAlpha(c) {
			s.scanIdentifier()
		} else {
			loxError(s.line, fmt.Sprintf("unexpected character: %c", c))
		}
	}

}

func (s *Scanner) ScanTokens() error {
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

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
	}
}
