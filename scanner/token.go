package scanner

type TokenType int

const (
	LeftCurlyBracket TokenType = iota
	RightCurlyBracket
	LeftSquareBracket
	RightSquareBracket
	String
	Number
	Boolean
	Colon
	Comma
	Null
	Eof
)

type Token struct {
	TokenType TokenType
	Value     any
}
