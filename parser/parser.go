package parser

import (
	"bufio"
	// "fmt"
	// "io"
	// "io/ioutil"
	"os"
)

// type Parserer interface {
// 	New(*os.File )
// 	HasMoreCommands()
// 	Advance()
// 	CommandType()
// 	Arg1()
// 	Arg2()
// }

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
