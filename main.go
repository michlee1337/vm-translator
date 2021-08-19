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

	file_name := os.Args[1][:strings.Index(os.Args[1], ".vm")]
	outfile, err := os.Create(fmt.Sprintf("%v.asm", file_name))
	if err != nil {
		panic("failed to create file")
	}
	defer outfile.Close()

	p := parser.New(infile)
	c := coder.New(file_name)

	for p.Advance() {
		switch p.CommandType() {
			case parser.CPush:
				outfile.WriteString(c.WritePush(p.Arg1(), p.Arg2()))
			case parser.CPop:
				outfile.WriteString(c.WritePop(p.Arg1(), p.Arg2()))
			case parser.CArithmetic:
				outfile.WriteString(c.WriteArithmetic(p.Arg1()));
		}
	}
}
