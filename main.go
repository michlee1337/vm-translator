package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/michlee1337/vm-translator/parser"
	"github.com/michlee1337/vm-translator/coder"
)

func main() {
	if len(os.Args) < 2 {
		panic("vm translator should be called with an argument as follows: vm-translator [file-to-translate.asm]")
	}

	infile, err := os.Open(os.Args[1])
	if err != nil {
		panic("failed to open file")
	}
	parser := parser.New(infile)
	parser.Advance()
	fmt.Println(parser.CurLine())
}
