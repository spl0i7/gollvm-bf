# gollvm-bf

![banner](res/image.jpg)

This is a toy project exploring viability of generating LLVM IR from brainfuck source in Go. It has two parts

(a) Interpreter -  Typical interpreter implementation of the brainfuck for referencing program output

(b) Compiler - Compiles source into LLVM IR which can then further be compiled to binary via clang.

