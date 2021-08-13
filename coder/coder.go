package coder

import (
	"fmt"
	"strings"
)

// stores segment pointers for writePushPop
// static segment is handled seperately
var segment_ptr = map[string]string{
	"local"			: "LCL",
	"argument"	: "ARG",
	"this"			: "THIS",
	"that"			: "THAT",
	"constant"	: "0",
	"temp"			: "5",
	"pointer"		: "3"}

type Coder struct {
	file_name string
}

func New(file_name string) *Coder {
	return &Coder{file_name}
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