package parser

import (
	"errors"
	"fmt"
	"jsonp/scanner"
)

type Parser struct {
	tokens []scanner.Token
	index  int
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{tokens: tokens, index: 0}
}

func (p *Parser) currentToken() scanner.Token {
	return p.tokens[p.index]
}

func (p *Parser) advance() {
	p.index++
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.index+1]
}

func (p *Parser) consume(token scanner.TokenType) error {
	currentToken := p.currentToken()
	if currentToken.TokenType == token {
		p.advance()
		return nil
	}
	return errors.New(fmt.Sprintf("Expected %v but found %v", token, currentToken))
}

func (p *Parser) Parse() (*Tree, error) {
	t, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	return t, p.consume(scanner.Eof)
}

func (p *Parser) parseObject() (*Tree, error) {
	tree := &Tree{NodeType: "Object"}
	err := p.consume(scanner.LeftCurlyBracket)
	if err != nil {
		return nil, err
	}
	if p.currentToken().TokenType == scanner.RightCurlyBracket {
		return tree, p.consume(scanner.RightCurlyBracket)
	}
	memberTree, err := p.parseMember()
	if err != nil {
		return nil, err
	}
	tree.Children = append(tree.Children, memberTree)
	for p.currentToken().TokenType == scanner.Comma {
		err := p.consume(scanner.Comma)
		if err != nil {
			return nil, err
		}
		memberTree, err := p.parseMember()
		if err != nil {
			return nil, err
		}
		tree.Children = append(tree.Children, memberTree)
	}
	return tree, p.consume(scanner.RightCurlyBracket)
}

func (p *Parser) parseMember() (*Tree, error) {
	tree := &Tree{NodeType: "Member"}
	stringToken := p.currentToken()
	err := p.consume(scanner.String)
	if err != nil {
		return nil, err
	}
	tree.Children = append(tree.Children, &Tree{Token: stringToken})
	colon := p.currentToken()
	err = p.consume(scanner.Colon)
	if err != nil {
		return nil, err
	}
	tree.Children = append(tree.Children, &Tree{Token: colon})
	valueTree, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	tree.Children = append(tree.Children, valueTree)
	return tree, nil
}

// TODO: Add support for null values
func (p *Parser) parseValue() (*Tree, error) {
	currentToken := p.currentToken()
	switch currentToken.TokenType {
	case scanner.Boolean, scanner.Number, scanner.String, scanner.Null:
		return &Tree{Token: currentToken}, p.consume(currentToken.TokenType)
	case scanner.LeftSquareBracket:
		return p.parseArray()
	case scanner.LeftCurlyBracket:
		return p.parseObject()
	default:
		return nil, errors.New("Unpexpected value in value field")
	}
}

func (p *Parser) parseArray() (*Tree, error) {
	tree := &Tree{NodeType: "Array"}
	err := p.consume(scanner.LeftSquareBracket)
	if err != nil {
		return nil, err
	}
	if p.currentToken().TokenType == scanner.RightSquareBracket {
		return tree, p.consume(scanner.RightSquareBracket)
	}
	valueTree, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	tree.Children = append(tree.Children, valueTree)
	for p.currentToken().TokenType == scanner.Comma {
		err := p.consume(scanner.Comma)
		if err != nil {
			return nil, err
		}
		valueTree, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		tree.Children = append(tree.Children, valueTree)
	}
	return tree, p.consume(scanner.RightSquareBracket)
}
