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

	file_name := os.Args[1][:strings.Index(os.Args[1], ".")]
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
				arg1, _ := p.Arg1()
				arg2, _ := p.Arg2()
				outfile.WriteString(c.WritePush(arg1, arg2))
			case parser.CPop:
				arg1, _ := p.Arg1()
				arg2, _ := p.Arg2()
				outfile.WriteString(c.WritePop(arg1, arg2))
			case parser.CArithmetic:
				arg1, _ := p.Arg1()
				outfile.WriteString(c.WriteArithmetic(arg1));

		}
	}
}
