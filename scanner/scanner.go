package scanner

import (
	"errors"
	"fmt"
)

const (
	boolFalse = "false"
	boolTrue  = "true"
)

type Scanner struct {
	source   []byte
	position int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source:   []byte(source),
		position: 0,
	}
}

func (s *Scanner) Scan() ([]Token, error) {
	tokens := []Token{}
	for s.position < len(s.source) {
		currentChar := s.getCurrentChar()
		switch *currentChar {
		case ' ', '\n', '\t':
			s.advance()
		case '{':
			tokens = append(tokens, Token{TokenType: LeftCurlyBracket, Value: "{"})
			s.advance()
		case '}':
			tokens = append(tokens, Token{TokenType: RightCurlyBracket, Value: "}"})
			s.advance()
		case '[':
			tokens = append(tokens, Token{TokenType: LeftSquareBracket, Value: "["})
			s.advance()
		case ']':
			tokens = append(tokens, Token{TokenType: RightSquareBracket, Value: "]"})
			s.advance()
		case ':':
			tokens = append(tokens, Token{TokenType: Colon, Value: ":"})
			s.advance()
		case ',':
			tokens = append(tokens, Token{TokenType: Comma, Value: ","})
			s.advance()
		case 'f', 't':
			token, err := s.parseBoolean()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *token)
		case 'n':
			token, err := s.parseNull()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *token)
		case '"':
			// parse string
			token, err := s.parseString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *token)
		// TODO: Add support for negative and float numbers
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			token, err := s.parseNumber()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *token)
		default:
			errorMsg := fmt.Sprintf("No action Found for %s", string(*currentChar))
			return nil, errors.New(errorMsg)
		}
	}
	tokens = append(tokens, Token{TokenType: Eof})
	return tokens, nil
}

func (s *Scanner) parseNull() (*Token, error) {
	return s.parseLiteral("null")
}

func (s *Scanner) parseNumber() (*Token, error) {
	number := 0
	for s.getCurrentChar() != nil && '0' <= *s.getCurrentChar() && *s.getCurrentChar() <= '9' {
		number *= 10
		number += int(*s.getCurrentChar() - '0')
		s.advance()
	}
	return &Token{TokenType: Number, Value: number}, nil
}

func (s *Scanner) getCurrentChar() *byte {
	if s.position < len(s.source) {
		return &s.source[s.position]
	}
	return nil
}

func (s *Scanner) consume(expected byte) error {
	currentChar := s.getCurrentChar()
	if currentChar == nil {
		errorMsg := fmt.Sprintf("End of file but expected %s", string(expected))
		return errors.New(errorMsg)
	}
	if *currentChar != expected {
		errorMsg := fmt.Sprintf("Expected %s but found %s\n", string(expected), string(*s.getCurrentChar()))
		return errors.New(errorMsg)
	}
	s.position++
	return nil
}

func (s *Scanner) advance() {
	s.position++
}

func (s *Scanner) parseString() (*Token, error) {
	err := s.consume('"')
	if err != nil {
		return nil, err
	}
	parsedString := []byte{}
	for s.getCurrentChar() != nil && *s.getCurrentChar() != '"' {
		parsedString = append(parsedString, *s.getCurrentChar())
		s.advance()
	}
	err = s.consume('"')
	if err != nil {
		return nil, err
	}
	stringToken := Token{TokenType: String, Value: string(parsedString)}
	return &stringToken, nil
}

func (s *Scanner) parseBoolean() (*Token, error) {
	currentChar := s.getCurrentChar()
	if *currentChar == 'f' {
		return s.parseBoolLiteral(boolFalse)
	}
	return s.parseBoolLiteral(boolTrue)
}

func (s *Scanner) parseBoolLiteral(lit string) (*Token, error) {
	t, err := s.parseLiteral(lit)
	if err != nil {
		return nil, err
	}
	t.Value = lit == "true"
	return t, nil
}

func (s *Scanner) parseLiteral(lit string) (*Token, error) {
	for _, v := range lit {
		currentChar := s.getCurrentChar()
		if currentChar == nil || *currentChar != byte(v) {
			return nil, errors.New("Unexpected char in value")
		}
		s.advance()
	}
	return &Token{TokenType: Boolean, Value: lit}, nil
}
