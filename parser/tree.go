package parser

import (
	"fmt"
	"jsonp/scanner"
	"strings"
)

type Tree struct {
	Token    scanner.Token
	NodeType string
	Children []*Tree
}

func (t Tree) String() string {
	return t.print("")
}

func (t Tree) print(prefix string) string {
	var s strings.Builder
	if len(t.NodeType) > 0 {
		s.WriteString(prefix)
		fmt.Fprintf(&s, "NodeType - %s\n", t.NodeType)
	}
	if t.Token.Value != nil {
		s.WriteString(prefix)
		fmt.Fprintf(&s, "Value: %+v\n", t.Token.Value)
	}
	for _, tt := range t.Children {
		s.WriteString(tt.print(prefix + "    "))
	}
	return s.String()
}
