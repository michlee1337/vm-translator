# vm-translator
Project 07 from Nand2Tetris Coursera course part 2.

The executable takes in a file written in the course's VM language and creates a new file with a .asm language containing the equivalent program in assembly language.

To use, run the following in a terminal
`./vm-translator [path-to-file/file-to-be-translated.vm]`

It contains three main packages:
- parser: responsible for parsing commands and command arguments from the input file
- coder: responsible for translating vm commands into assembly commands
- main: handles file io, passes parsed output to the coder, and writes coder output to the .asm file
