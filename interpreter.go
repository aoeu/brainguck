package brainguck

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

type Interpreter struct {
	out        *os.File
	index      int
	offset     int
	cells      []byte
	blockStack []int
	storeNext  bool
	print      bool
	skipBlock  int
}

func InterpretFile(filename string) (n int, err error) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	i := NewInterpreter()
	return i.Read(input)
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		out:        os.Stdout,
		index:      0,
		cells:      make([]byte, 30000),
		blockStack: make([]int, 0),
		skipBlock:  -1,
	}
}

func (i *Interpreter) Read(in []byte) (n int, err error) {
	var b byte
	stdin := bufio.NewReader(os.Stdin)
	for i.offset = 0; i.offset != len(in); n++ {
		b = i.interpret(in[i.offset])
		if i.print {
			fmt.Fprintf(i.out, "%c", b)
			i.print = false
		}
		if i.storeNext {
			b, err = stdin.ReadByte()
			if err != nil {
				return n, err
			}
			i.interpret(b)
		}

	}
	return n, err
}

// TODO(aoeu): Satisfy Writer interface.

/*
Character	Meaning
	>		increment the data pointer (to point to the next cell to the right).
	<		decrement the data pointer (to point to the next cell to the left).
	+		increment (increase by one) the byte at the data pointer.
	-		decrement (decrease by one) the byte at the data pointer.
	.		output the byte at the data pointer.
	,		accept one byte of input, storing its value in the byte at the data pointer.
	[		if the byte at the data pointer is zero, then instead of moving the instruction
				pointer forward to the next command, jump it forward to the command after the matching ] command.
	]		if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward
				to the next command, jump it back to the command after the matching [ command.
*/

func (i *Interpreter) interpret(b byte) byte {
	if i.storeNext {
		i.cells[i.index], i.storeNext = b, false
		return 0
	}
	i.offset++
	// Handle loop control operators (and looping conditions) first.
	switch {
	case b == ']':
		n := len(i.blockStack)
		var popped int
		popped, i.blockStack = i.blockStack[n-1], i.blockStack[:n-1]
		switch i.skipBlock {
		case -1:
			if i.cells[i.index] != 0 {
				i.offset = popped
			}
		case n:
			i.skipBlock = -1
		}
		return 0
	case b == '[':
		i.blockStack = append(i.blockStack, i.offset-1)
		if i.skipBlock == -1 {
			if i.cells[i.index] == 0 {
				i.skipBlock = len(i.blockStack)
			}
		}
	case i.skipBlock > 0:
		return 0
	}
	// Handle other operators.
	switch b {
	case '>':
		i.index++
	case '<':
		if i.index > 0 {
			i.index--
		}
	case '+':
		i.cells[i.index]++
	case '-':
		i.cells[i.index]--
	case '.':
		i.print = true
		return i.cells[i.index]
	case ',':
		i.storeNext = true
	}
	return 0
}
