package main

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"github.com/michlee1337/vm-translator/parser"
	"github.com/michlee1337/vm-translator/coder"
)

func main() {
	debugPtr := flag.Bool("debug", false, "Generates .asm files that are more readable for debugging")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		panic("vm translator should be called with an argument as follows: vm-translator [-debug] [file-to-translate.asm]")
	}

	infile, err := os.Open(args[0])
	if err != nil {
		panic("failed to open file")
	}

	file := args[0][:strings.Index(args[0], ".vm")]
	outfile, err := os.Create(fmt.Sprintf("%v.asm", file))
	if err != nil {
		panic("failed to create file")
	}
	defer outfile.Close()

	p := parser.New(infile)
	c := coder.New(file[strings.LastIndex(file, "/")+1:], *debugPtr)
	
	if *debugPtr {
		outfile.WriteString("// This file was generated in debug mode\n")
	}
	for p.Advance() {
		switch p.CommandType() {
			case parser.CPush:
				outfile.WriteString(c.WritePush(p.Arg1(), p.Arg2()))
			case parser.CPop:
				outfile.WriteString(c.WritePop(p.Arg1(), p.Arg2()))
			case parser.CArithmetic:
				outfile.WriteString(c.WriteArithmetic(p.Arg1()))
			case parser.CLabel:
				outfile.WriteString(c.WriteLabel(p.Arg1()))
			case parser.CGoto:
				outfile.WriteString(c.WriteGoto(p.Arg1()))
			case parser.CIf:
				outfile.WriteString(c.WriteIf(p.Arg1()))
		}
	}

	outfile.WriteString(c.WriteClose());
	fmt.Println("Translated successfully! New file written at ", outfile.Name())
}
