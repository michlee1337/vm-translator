// This is a test file for package parser
// To run tests, run go test -v
package coder

import (
	"fmt"
	"testing"
)

// func testCoder() {
// 	c := New("test", false)
// 	m.Run()
// }

// func TestPushConstant(t *testing.T) {
// 	c := New("test", false)
// 	code := c.WritePush("constant", "9")
// 	expected := "@9\n" +
// 							"D=A\n" +
// 							"@SP\n" +
// 							"A=M\n" +
// 							"M=D\n" + 
// 							"@SP\n" +
// 							"M=M+1\n"
// 	if code != expected {
// 		t.Errorf("WritePush(\"constant\", \"9\")) = %s; want %s", code, expected)
// 	}
// }

func TestPushConstant(t *testing.T) {
	var tests = []struct {
		segment, addr string
		code string
	}{
			{"constant", "10", 
				"@10\n" +
				"D=A\n" +
				"@SP\n" +
				"A=M\n" +
				"M=D\n" +
				"@SP\n" +
				"M=M+1\n"},
			{"constant", "45",
				"@45\n" +
				"D=A\n" +
				"@SP\n" +
				"A=M\n" +
				"M=D\n" +
				"@SP\n" +
				"M=M+1\n"},
			{"constant", "510",
				"@510\n" +
				"D=A\n" +
				"@SP\n" +
				"A=M\n" +
				"M=D\n" +
				"@SP\n" +
				"M=M+1\n"},
		}
		c := New("test", false)

	for _, tt := range tests {
			testname := fmt.Sprintf("%s,%s", tt.segment, tt.addr)
			t.Run(testname, func(t *testing.T) {
					code := c.WritePush(tt.segment, tt.addr)
					if code != tt.code {
							t.Errorf("got %s, want %s", code, tt.code)
					}
			})
	}
}