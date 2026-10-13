package main

import (
	"fmt"
	"jsonp/parser"
	"jsonp/scanner"
)

func main() {
	inputString := `{
        "Image": {
            "Width":  800,
            "Height": 600,
            "Title":  "View from 15th Floor",
            "Thumbnail": {
                "Url":    "http://www.example.com/image/481989943",
                "Height": 125,
                "Width":  100
            },
            "Animated" : false,
			"ExtraData": null,
            "IDs": [116, 943, 234, 38793]
          }
      }`
	scanner := scanner.NewScanner(inputString)
	tokens, err := scanner.Scan()
	if err != nil {
		fmt.Println(err)
		return
	}
	parser := parser.NewParser(tokens)
	tree, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tree)
}
