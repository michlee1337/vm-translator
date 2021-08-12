package parser

import (
	"bufio"
	// "fmt"
	// "io"
	// "io/ioutil"
	"os"
	"strings"
)

type CommandType int

const (
	CArithmetic CommandType = iota
  CPush
  CPop
  // CGoto  // TODO: Support these command types
	// CIf
	// CFunction
	// CReturn
	// CCall
)

type Parser struct {
	in    *os.File
	scanner *bufio.Scanner
}

func New(in *os.File) *Parser {
	reader := bufio.NewReader(in)
	scanner := bufio.NewScanner(reader)
	return &Parser{in, scanner}
}

func (p *Parser) HasMoreLines() bool {
	return p.scanner.Scan()
}

func (p *Parser) Advance() {
	if !p.HasMoreLines() {
		panic("No more lines")
	}
}

func (p *Parser) CurLine() (string) {
	return p.scanner.Text()
}

func (p *Parser) CommandType() (CommandType) {
	first_word := strings.Fields(p.CurLine())[0]
	if first_word == "push" {
		return CPush
	} else if first_word == "pop" {
		return CPop
	} else {  // TODO: Handle other commands
		return CArithmetic
	}
}

func (p *Parser) Arg1() (string) {}

func (p *Parser) Arg2() (string) {}
