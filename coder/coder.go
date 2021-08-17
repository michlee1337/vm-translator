// Package coder implements functions for translating
// commands in VM language to Assembly language.
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

// Map of VM comparators to their negated Assembly jump commands.
// This is used to implement branching to translate comparator 
// arithmetic commands.
var cmp_false = map[string]string{
	"GT": "JLE",  // greater than == not less than or equal to
	"LT": "JGE",  // less than == not greater than or equal to
	"EQ": "JNE"}  // equal == not not equal
		
// Map of number of times this particular VM comparator has
// appeared in the current instream being translated.
// This is used to create unique branch labels to translate
// comparator arithmetic commands.
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
	file_name string  // necessary for handling static segment commands
}

func New(file_name string) *Coder {
	return &Coder{file_name}
}

// Translates commands of type CPush
func (c *Coder) WritePush(segment string, addr string) string {
	var sb strings.Builder

	// move to segment
	sb.WriteString(c.getSegment(segment, addr))

	// copy value from segment into stack
	sb.WriteString(
		"D=M\n" +
		"@SP\n" +
		"M=D\n")

	// increment stack pointer
	sb.WriteString(
		"@SP\n" +
		"M=M+1\n")

	return sb.String()
}

// Translates commands of type CPop
func (c *Coder) WritePop(segment string, addr string) string {
	var sb strings.Builder

	// move to stack and store data
	sb.WriteString(
		"@SP\n" +
		"A=M\n" +
		"D=M\n")

	// move to segment
	sb.WriteString(c.getSegment(segment, addr))

	// copy data into segment
	sb.WriteString(
		"M=D\n")

	// decrement stack pointer
	sb.WriteString(
		"@SP\n" +
		"M=M-1\n")

	return sb.String()
}

// Translates commands of type CArithmetic
func (c *Coder) WriteArithmetic(op string, err error) string{
	switch op {
		case "add":
			return 	goto_topmost_stack_val +
							pop_into_D +
							"M=M+D\n" +
							decrement_SP, nil
		case "sub":
			return 	goto_topmost_stack_val +
							pop_into_D +
							"M=M-D\n" +
							decrement_SP, nil
		case "gt":
			return 	goto_topmost_stack_val +
							pop_into_D +
							c.writeCompResultToStack("gt"), nil
		case "lt":
			return 	goto_topmost_stack_val +
							pop_into_D +
							c.writeCompResultToStack("lt"), nil
		case "eq":
			return 	goto_topmost_stack_val +
							pop_into_D +
							c.writeCompResultToStack("eq"), nil
		case "neg":
			return 	 +
							"M=-M", nil
		case "and":
			return 	goto_topmost_stack_val +
							pop_into_D +
							"M=M&D\n" +
							decrement_SP, nil
		case "or":
			return 	goto_topmost_stack_val +
							pop_into_D +
							"M=M|D\n" +
							decrement_SP, nil
		case "not":
			return 	goto_topmost_stack_val +
							"M=!M", nil
	}
	return "", errors.New("Command is not valid")
}

// Moves memory pointer to the specified segment address
func (c *Coder) getSegment(segment string, addr string) string {
	// handle static as special case
	if segment == "static" {
		return fmt.Sprintf("@%v.%v\n", c.file_name, addr)
	}

	return "@" + segment_ptr[segment] + "\n" +
					"D=A\n" +
					"@" + addr + "\n" +
					"A=D+A\n"
}

// Writes boolean result of comparison to the stack
// true is respresented as -1 (0b111...111), false as 0 (0b000...000)
func (c *Coder) writeCompResultToStack(comparator string) string {
	label := fmt.Sprintf("NOT_%v%d", comparator, label_count[comparator])
	cond := cmp_false[comparator]
	return	"D=M-D\n" + 					// D = second topmost val - topmost val
					"@SP\n" +  						
					"M=M+1" + 						// preemptively increment stack pointer 
					"A=M-1\n" +							// goto top of stack

					"@" + label + "\n" +  // jump if comparator result is false
					"D;" + cond +"\n" +

					"M=-1\n" +  // if true, write -1 and end
					"@END\n" +
					"0; JMP\n" +

					"(" + label + ")\n" +  // if false, write 0 and end
					"M=0\n" +

					"(END)\n" +  // loop infinitely
					"@END\n" +
					"0; JMP\n"
} 