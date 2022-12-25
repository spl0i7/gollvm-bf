package main

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/marstr/collection/v2"
	"io"
	"log"
	"os"
)

const (
	ErrUnexpectedBracket = "unexpected number of brackets"
	MemSize              = 30_000
)

type Loop struct {
	Body *ir.Block
	End  *ir.Block
}

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	program, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	mod := ir.NewModule()

	puts := mod.NewFunc("putchar", types.I32, ir.NewParam("", types.I32))
	memset := mod.NewFunc("memset", types.Void, ir.NewParam("p1", types.I8Ptr), ir.NewParam("p2", types.I8), ir.NewParam("p3", types.I64))

	entryPoint := mod.NewFunc("main", types.I32)

	builder := entryPoint.NewBlock("")

	st := collection.NewStack[Loop]()

	arrayType := types.NewArray(MemSize, types.I8)
	cellMemory := builder.NewAlloca(arrayType)

	builder.NewCall(memset,
		builder.NewGetElementPtr(arrayType, cellMemory, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0)),
		constant.NewInt(types.I8, 0),
		constant.NewInt(types.I64, MemSize),
	)

	dataPtr := builder.NewAlloca(types.I64)
	builder.NewStore(constant.NewInt(types.I64, 0), dataPtr)

	for i := 0; i < len(program); i++ {
		switch program[i] {

		case '+':
			ptr := builder.NewGetElementPtr(arrayType, cellMemory, constant.NewInt(types.I64, 0), builder.NewLoad(types.I64, dataPtr))
			added := builder.NewAdd(builder.NewLoad(types.I8, ptr), constant.NewInt(types.I8, 1))
			builder.NewStore(added, ptr)
		case '-':
			ptr := builder.NewGetElementPtr(arrayType, cellMemory, constant.NewInt(types.I64, 0), builder.NewLoad(types.I64, dataPtr))
			added := builder.NewAdd(builder.NewLoad(types.I8, ptr), constant.NewInt(types.I8, -1))
			builder.NewStore(added, ptr)
		case '<':
			t1 := builder.NewAdd(builder.NewLoad(types.I64, dataPtr), constant.NewInt(types.I8, -1))
			builder.NewStore(t1, dataPtr)
		case '>':
			t1 := builder.NewAdd(builder.NewLoad(types.I64, dataPtr), constant.NewInt(types.I8, +1))
			builder.NewStore(t1, dataPtr)
		case '.':
			ptr := builder.NewGetElementPtr(arrayType, cellMemory, constant.NewInt(types.I64, 0), builder.NewLoad(types.I64, dataPtr))
			builder.NewCall(puts, builder.NewSExt(builder.NewLoad(types.I8, ptr), types.I32))
		case '[':
			ptr := builder.NewGetElementPtr(arrayType, cellMemory, constant.NewInt(types.I64, 0), builder.NewLoad(types.I64, dataPtr))
			ld := builder.NewLoad(types.I8, ptr)

			cmpResult := builder.NewICmp(enum.IPredNE, ld, constant.NewInt(types.I8, 0))

			wb := Loop{
				Body: entryPoint.NewBlock(""),
				End:  entryPoint.NewBlock(""),
			}
			st.Push(wb)

			builder.NewCondBr(cmpResult, wb.Body, wb.End)
			builder = wb.Body

		case ']':

			front, ok := st.Pop()
			if !ok {
				log.Fatalln(ErrUnexpectedBracket)
			}

			ptr := builder.NewGetElementPtr(arrayType, cellMemory, constant.NewInt(types.I64, 0), builder.NewLoad(types.I64, dataPtr))
			ld := builder.NewLoad(types.I8, ptr)

			cmpResult := builder.NewICmp(enum.IPredNE, ld, constant.NewInt(types.I8, 0))

			builder.NewCondBr(cmpResult, front.Body, front.End)
			builder = front.End
		}
	}

	builder.NewRet(constant.NewInt(types.I32, 0))

	fmt.Println(mod.String())

}
