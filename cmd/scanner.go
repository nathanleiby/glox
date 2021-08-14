package cmd

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   string // TODO: Object?
	line      int
}

func (t Token) String() string {
	return fmt.Sprintf("%03d %s %s", t.tokenType, t.lexeme, t.literal)
}

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

// match is a "conditional advance()".
// It checks if the current character matches, and advances if so/
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	// TODO: Is this OK refactor?
	// if s.peek() != expected {
	// 	return false
	// }

	s.current = s.current + 1
	return true
}

func (s *Scanner) scanStringToken() {
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
	// value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING) //  TODO: addTokenliteral with real value
}

func (s *Scanner) addToken(tt TokenType) {
	text := s.source[s.start:s.current]
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
		s.scanStringToken()
	default:
		loxError(s.line, fmt.Sprintf("unexpected character: %c", c))
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
