// Package parser implements functions for parsing commands in the VM Language
// and an iota of command types in the VM Language.

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
	scanner *bufio.Scanner
	cur_line string
}

// Initializes new Parser struct with
// a bufio.Scanner for the instream
func New(in *os.File) *Parser {
	reader := bufio.NewReader(in)
	scanner := bufio.NewScanner(reader)
	return &Parser{scanner, ""}
}

// Sets Parser.cur_line to the next command and returns true.
// If no more commands, reryrbs false.
func (p *Parser) Advance() bool {
	for p.scanner.Scan() {
		p.cur_line = p.scanner.Text()
		strings.TrimSpace(p.cur_line)
    if p.isCommand() {
			return true
		}
	}
	return false
}

func (p *Parser) CurLine() (string) {
	return p.cur_line
}

// Returns an iota representing the type of VM command.
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

// Returns the first argument of the current command.
// If no such argument exists, returns error.
func (p *Parser) Arg1() (string, error) {
	if p.cur_line == "" {
		return "", errors.New("Parser has not yet started processing the file. Call parser.Advance() to start processing.")
	}
	cmd_type := p.CommandType()
	if cmd_type == CReturn {
		return "", errors.New("Commands of type Return has no arguments")
	}
	if cmd_type == CArithmetic {
		return strings.Fields(p.cur_line)[0], nil
	}
	return strings.Fields(p.cur_line)[1], nil
}


// Returns the second argument of the current command.
// If no such argument exists, returns error.
func (p *Parser) Arg2() (string, error) {
	if p.cur_line == "" {
		return "", errors.New("Parser has not yet started processing the file. Call parser.Advance() to start processing.")
	}
	cmd_type := p.CommandType()

	if cmd_type != CPush &&
			cmd_type != CPop &&
			cmd_type != CFunction &&
			cmd_type != CCall {
		return "", fmt.Errorf("Command of type %v has no second argument.", cmd_type)
	}
	return strings.Fields(p.cur_line)[2], nil
}

// Returns true if the current command is not whitespace and not a comment.
func (p *Parser) isCommand() bool {
	return len(p.cur_line) > 1 && p.cur_line[:1] != "//"
}
