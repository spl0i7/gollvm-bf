package main

import (
	"errors"
	"fmt"
	"github.com/marstr/collection/v2"
	"io"
	"log"
	"os"
)

const (
	Increase  = '+'
	Decrease  = '-'
	MoveRight = '>'
	MoveLeft  = '<'
	Output    = '.'
	Input     = ','
	JumpRight = '['
	JumpLeft  = ']'
)

type Interpreter struct {
	source  []byte
	memory  []byte
	pointer int
	pc      int
	jumpMap map[int]int
}

func (i *Interpreter) Run() error {

	s := collection.NewStack[int]()

	for i.pc = 0; i.pc < len(i.source); i.pc++ {
		switch i.source[i.pc] {
		case JumpRight:
			s.Push(i.pc)
		case JumpLeft:
			top, ok := s.Pop()
			if !ok {
				return errors.New("invalid program")
			}
			i.jumpMap[top] = i.pc
			i.jumpMap[i.pc] = top
		}
	}

	for i.pc = 0; i.pc < len(i.source); i.pc++ {
		switch i.source[i.pc] {
		case Increase:
			i.memory[i.pointer]++
		case Decrease:
			i.memory[i.pointer]--
		case MoveRight:
			i.pointer++
			i.pointer %= len(i.memory)
		case MoveLeft:
			i.pointer--
			i.pointer += len(i.memory)
			i.pointer %= len(i.memory)
		case Output:
			fmt.Print(string(i.memory[i.pointer]))
		case Input:
			if _, err := fmt.Scanf("%c", &i.memory[i.pointer]); err != nil {
				return err
			}
		case JumpRight:
			if i.memory[i.pointer] == 0 {
				i.pc = i.jumpMap[i.pc]
			}
		case JumpLeft:
			if i.memory[i.pointer] != 0 {
				i.pc = i.jumpMap[i.pc]
			}
		}
	}

	return nil

}

func NewInterpreter(source []byte) *Interpreter {
	return &Interpreter{
		source:  source,
		memory:  make([]byte, 30_000),
		jumpMap: map[int]int{},
	}
}

func main() {

	source, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
	}
	sourceCode, err := io.ReadAll(source)
	if err != nil {
		log.Println(err)
	}
	i := NewInterpreter(sourceCode)
	if err := i.Run(); err != nil {
		log.Println(err)
	}

}
