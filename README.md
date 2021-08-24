# vm-translator
This implements the two-tier compiler project (high-level -> virtual machine code -> machine language) from Coursera course Nand2Tetris part 2.

The executable takes in a file written in the course's VM language and creates a new file with a .asm language containing the equivalent program in assembly language.

## Instructions
Run the following in a terminal
`./vm-translator [path-to-file/file-to-be-translated.vm]`
A .asm file with the same name will be generated in the same folder.

## Code Organization
It contains three main packages:
- parser: responsible for parsing commands and command arguments from the input file
- coder: responsible for translating vm commands into assembly commands
- main: handles file io, passes parsed output to the coder, and writes coder output to the .asm file
