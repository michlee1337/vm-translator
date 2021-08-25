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
	"temp"			: "R5",
	"pointer"		: "R3"}

// Map of VM comparators to their negated Assembly jump commands.
// This is used to implement branching to translate comparator 
// arithmetic commands.
var cmp_false = map[string]string{
	"gt": "JLE",  // greater than == not less than or equal to
	"lt": "JGE",  // less than == not greater than or equal to
	"eq": "JNE"}  // equal == not not equal
		
// Map of number of times this particular VM comparator has
// appeared in the current instream being translated.
// This is used to create unique branch labels to translate
// comparator arithmetic commands.
var label_count = map[string]int{
	"gt": 0,
	"lt": 0,
	"eq": 0}

// Common sub-operation
const goto_topmost_stack_val = 	"@SP\n" +
															"A=M-1\n"

// Common sub-operation
const pop_into_D = 	"D=M\n" + 	 	// store val in D
										"A=A-1\n"   	// move to second top-most

// Common sub-operation
const decrement_SP = 	"@SP\n" +
											"M=M-1\n"

// Coder implements translation from VM Language to Assembly Language.
type Coder struct {
	file_name string  // necessary for handling static segment commands
	debug bool 				// if true, generates .asm code that is more readable
}

func New(file_name string, debug bool) *Coder {
	return &Coder{file_name, debug}
}

// Translates commands of type CPush
func (c *Coder) WritePush(segment string, addr string) string {
	var sb strings.Builder

	if c.debug {
		sb.WriteString(
			"\n// push " + segment + addr + "\n")
	}

	if segment == "constant" {
		// constant is not a segment in storage space.
		// Instead it refers to the actual integer value.
		sb.WriteString(
			"@" + addr + "\n" +
			"D=A\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n")
	} else {
		// move to segment
		sb.WriteString(c.getSegment(segment, addr))

		// copy value from segment into stack
		sb.WriteString(
			"D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n")
	}

	// increment stack pointer
	sb.WriteString(
		"@SP\n" +
		"M=M+1\n")
	
	return sb.String()
}

// Translates commands of type CPop
func (c *Coder) WritePop(segment string, addr string) string {
	var sb strings.Builder

	if c.debug {
		sb.WriteString(
			"\n// pop " + segment + addr + "\n")
	}

	// save segment location
	sb.WriteString(
		c.getSegment(segment, addr) +
		"D=A\n" +
		"@R13\n" +
		"M=D\n")

	// move to stack
	sb.WriteString(goto_topmost_stack_val)
	
	// copy data into segment
	sb.WriteString(
		"D=M\n" +
		"@R13\n" +
		"A=M\n" +
		"M=D\n")

	// decrement stack pointer
	sb.WriteString(decrement_SP)

	return sb.String()
}

// Translates commands of type CArithmetic
func (c *Coder) WriteArithmetic(op string) string {
	var sb strings.Builder

	if c.debug {
		sb.WriteString(
			"\n// " + op + "\n")
	}
	switch op {
		case "add":
			sb.WriteString(
				goto_topmost_stack_val +
				pop_into_D +
				"M=M+D\n" +
				decrement_SP)
		case "sub":
			sb.WriteString(
				goto_topmost_stack_val +
				pop_into_D +
				"M=M-D\n" +
				decrement_SP)
		case "gt":
			sb.WriteString(
				goto_topmost_stack_val +
				pop_into_D +
				c.writeCompResultToStack("gt") +
				decrement_SP)
		case "lt":
			sb.WriteString(
				goto_topmost_stack_val +
				pop_into_D +
				c.writeCompResultToStack("lt") +
				decrement_SP)
		case "eq":
			sb.WriteString(
				goto_topmost_stack_val +
				pop_into_D +
				c.writeCompResultToStack("eq") +
				decrement_SP)
		case "neg":
			sb.WriteString(
				goto_topmost_stack_val +
				"M=-M\n")
		case "and":
			sb.WriteString(
				goto_topmost_stack_val +
				pop_into_D +
				"M=M&D\n" +
				decrement_SP)
		case "or":
			sb.WriteString(
				goto_topmost_stack_val +
				pop_into_D +
				"M=M|D\n" +
				decrement_SP)
		case "not":
			sb.WriteString(
				goto_topmost_stack_val +
				"M=!M\n")
		default:
			panic(fmt.Sprintf("Command %s is not valid", op))
	}
	return sb.String()
}

// Writes end of ASM files
func (c *Coder) WriteClose() string {
	return 	"(END)\n" +
					"@END\n" +
					"0; JMP\n"
}


// Moves memory pointer to the specified segment address
func (c *Coder) getSegment(segment string, addr string) string {
	// handle static as special case
	if segment == "static" {
		return fmt.Sprintf("@%v.%v\n", c.file_name, addr)
	}
	var sb strings.Builder

	// D = addr
	sb.WriteString(
		"@" + addr + "\n" +
		"D=A\n")

	// increment by segment start location
	if segment == "temp" || segment == "pointer" {
		sb.WriteString(
			"@" + segment_ptr[segment] + "\n" +
			"A=D+A\n")
	} else {
		sb.WriteString(
			"@" + segment_ptr[segment] + "\n" +
			"A=D+M\n")
	}

	return 	sb.String()
}

// Writes boolean result of comparison to the stack
// true is respresented as -1 (0b111...111), false as 0 (0b000...000)
func (c *Coder) writeCompResultToStack(comparator string) string {
	jump := fmt.Sprintf("NOT_%v%d", comparator, label_count[comparator])
	end := fmt.Sprintf("END_%v%d", comparator, label_count[comparator])
	label_count[comparator] += 1  // increment counter to ensure label uniqueness
	cond := cmp_false[comparator]
	return	"D=M-D\n" + 						// D = second topmost val - topmost val

					"@" + jump + "\n" +  		// jump if comparator result is false
					"D;" + cond +"\n" +

					"@SP\n" +  							// if true, write -1 and end
					"A=M-1\n" +
					"A=A-1\n" +
					"M=-1\n" +
					"@" + end + "\n" +
					"0; JMP\n" +

					"(" + jump + ")\n" +  	// if false, write 0 and end
					"@SP\n" +
					"A=M-1\n" +
					"A=A-1\n" +
					"M=0\n" +

					"(" + end + ")\n"
} 