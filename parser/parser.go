package parser

import (
	"bufio"
	// "fmt"
	// "io"
	// "io/ioutil"
	"os"
)

type CommandType int

const (
	Arithmetic CommandType = iota
  Push
  Pop
  // Goto  // TODO: Support these command types
	// If
	// Function
	// Return
	// Call
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

func (p *Parser) CommandType() (CommandType) {}

func (p *Parser) Arg1() (string) {}

func (p *Parser) Arg2() (string) {}
