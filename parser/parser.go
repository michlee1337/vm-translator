package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"errors"
)

type CommandType int

const (
	CArithmetic CommandType = iota
  CPush
  CPop
  CGoto
	CIf
	CFunction
	CReturn
	CCall
)

type Parser struct {
	in    *os.File
	scanner *bufio.Scanner
	cur_line string
}

func New(in *os.File) *Parser {
	reader := bufio.NewReader(in)
	scanner := bufio.NewScanner(reader)
	return &Parser{in, scanner, ""}
}

func (p *Parser) Advance() bool {
	// Returns false if no more lines
	for p.scanner.Scan() {
		p.cur_line = p.scanner.Text()
		strings.TrimSpace(p.cur_line)
    if p.IsCommand() {
			return true
		}
	}
	return false
}

func (p *Parser) CurLine() (string) {
	return p.cur_line
}

func (p *Parser) CommandType() (CommandType) {
	first_word := strings.Fields(p.cur_line)[0]
	if first_word == "push" {
		return CPush
	} else if first_word == "pop" {
		return CPop
	} else {  // TODO: Handle other commands
		return CArithmetic
	}
}

func (p *Parser) Arg1() (string) {}

func (p *Parser) IsCommand() bool {
	return len(p.cur_line) > 1 && p.cur_line[:1] != "//"
}
