package main

import (
	"fmt"
	"github.com/michlee1337/vm-translator/parser"
	"os"
)

func main() {
	fmt.Println("Hello")
	infile, err := os.Open(".gitignore")
	if err != nil {
		panic("failed to open file")
	}
	parser := parser.New(infile)
	parser.Advance()
	fmt.Println(parser.CurLine())
}
