package coder

import (
	"fmt"
	"strings"
)

// Map of segment pointers for the Hack machine spec.
// The static segment is handled in WritePushPop()
// as a special case as there is no fixed pointer.
var segment_ptr = map[string]string{
	"local"			: "LCL",
	"argument"	: "ARG",
	"this"			: "THIS",
	"that"			: "THAT",
	"constant"	: "0",
	"temp"			: "5",
	"pointer"		: "3"}

var false_condition = map[string]string{
	"GT": "JLE",
	"LT": "JGE",
	"EQ": "JNE"}
		
// track number of generated labels to ensure uniqueness
var label_count = map[string]int{
	"GT": 0,
	"LT": 0,
	"EQ": 0}

// Common sub-operation in arithmetic commands
const goto_topmost_stack_val = 	"@SP\n" +
															"A=M-1\n"

// Common sub-operation in arithmetic commands
const pop_into_D = 	"D=M\n" + 	 	// store val in D
										"A=A-1\n"   	// move to second top-most

// Common sub-operation in arithmetic commands
const decrement_SP = 	"@SP\n" +
											"M=M-1"

// Coder implements translation from VM Language to Assembly Language.
type Coder struct {
	file_name string
}

func New(file_name string) *Coder {
	return &Coder{file_name}
}

func (c *Coder) WritePush(segment string, addr string) string {
	var sb strings.Builder
	sb.WriteString(c.GetSegment(segment, addr))

	// move from segment to stack
	sb.WriteString(
		"D=M\n" +
		"@SP\n" +
		"M=D\n")

	// update stack pointer
	sb.WriteString(
		"@SP\n" +
		"M=M+1\n")

	return sb.String()
}

func (c *Coder) WritePop(segment string, addr string) string {
	var sb strings.Builder

	// get stack
	sb.WriteString(
		"@SP\n" +
		"A=M\n" +
		"D=M\n")

	// move from stack to segment
	sb.WriteString(c.GetSegment(segment, addr))
	sb.WriteString(
		"M=D\n")

	// update stack pointer
	sb.WriteString(
		"@SP\n" +
		"M=M-1\n")

	return sb.String()
}

func (c *Coder) WriteArithmetic(op string) string{
	switch op {
		case "add":
			return 	pop_into_D +
							"M=M+D\n" +
							decrement_SP
		case "sub":
			return 	pop_into_D +
							"M=M-D\n" +
							decrement_SP
		case "gt":
			return 	pop_into_D +
							c.ComparisonBranch("gt")
		case "lt":
			return 	pop_into_D +
							c.ComparisonBranch("lt")
		case "eq":
			return 	pop_into_D +
							c.ComparisonBranch("eq")
		// case "neg":
		// case "and":
		// case "or":
		// case "not":
	}
	return ""
}

func (c *Coder) GetSegment(segment string, addr string) string {
	// handle static as special case
	if segment == "static" {
		return fmt.Sprintf("@%v.%v\n", c.file_name, addr)
	}

	return "@" + segment_ptr[segment] + "\n" +
					"D=A\n" +
					"@" + addr + "\n" +
					"A=D+A\n"
}

func (c *Coder) ComparisonBranch(comparator string) string {
	jump_label := fmt.Sprintf("NOT_%v%d", comparator, label_count[comparator])
	jump_cond := false_condition[comparator]
	return	"D=M-D\n" + // D = second topmost (v1) - topmost (v2)
					"@SP\n" +  // goto store location
					"A=M\n" +

					"@" + jump_label + "\n" +  // jump if v1 gt v2 is false
					"D;" + jump_cond +"\n" +

					"M=-1\n" +  // if gt
					"@END\n" +
					"0; JMP\n" +

					"(" + jump_label + ")\n" +  // if not gt
					"M=0\n" +

					"(END)\n" +
					"@END\n" +
					"0; JMP\n"
} 