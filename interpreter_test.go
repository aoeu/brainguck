package brainguck

import (
	"testing"
)

/*
Character	Meaning
	>		increment the data pointer (to point to the next cell to the right).
	<		decrement the data pointer (to point to the next cell to the left).
	+		increment (increase by one) the byte at the data pointer.
	-		decrement (decrease by one) the byte at the data pointer.
	.		output the byte at the data pointer.
	,		accept one byte of input, storing its value in the byte at the data pointer.
	[			if the byte at the data pointer is zero, then instead of moving the instruction
			pointer forward to the next command, jump it forward to the command after the matching ] command.
	]		if the byte at the data pointer is nonzero, then instead of moving the instruction pointer forward
				to the next command, jump it back to the command after the matching [ command.
*/

func TestInterpretIndexOps(t *testing.T) {
	testCases := []struct {
		in    byte
		index int
	}{
		{'>', 1},
		{'<', 0},
	}
	for _, tc := range testCases {
		intr := NewInterpreter()
		intr.interpret(tc.in)
		if tc.index != intr.index {
			t.Errorf("Expected %v but resulted in %v\n", tc.index, intr.index)
		}
	}
}

func TestInterpretCellOps(t *testing.T) {
	testCases := []struct {
		in    byte
		index int
		data  byte
	}{
		{'+', 0, 1},
		{'-', 0, 255},
	}
	for _, tc := range testCases {
		intr := NewInterpreter()
		intr.interpret(tc.in)
		e, a := tc.data, intr.cells[intr.index]
		if e != a {
			t.Errorf("Expected %h but cell data was %h\n", e, a)
		}

	}
}

func TestInterpretOutputOp(t *testing.T) {
	intr := NewInterpreter()
	var expected byte = 'a'
	intr.cells[0] = expected
	actual := intr.interpret('.')
	if expected != actual {
		t.Errorf("Expected '%v' but received '%v'", expected, actual)
	}
}

func TestInterpretInputOp(t *testing.T) {
	// Input
	intr := NewInterpreter()
	intr.interpret(',')
	var expected byte = 'a'
	intr.interpret(expected)
	actual := intr.cells[intr.index]
	if expected != actual {
		t.Errorf("Expected '%v' but received '%v'", expected, actual)
	}
}

func TestSkipBlocks(t *testing.T) {
	intr := NewInterpreter()
	intr.interpret('[')
	if intr.skipBlock == -1 {
		t.Errorf("Interpreter should skip current block.")
	}
	for _, op := range []byte{'+', '-', '<', '>', '[', '.', '[', ']', ',', ']'} {
		b := intr.interpret(op)
		if b != 0 || intr.index != 0 || intr.cells[intr.index] != 0 {
			t.Errorf("Interpreter should be ignoring operators until end of current block.")
		}
		if intr.skipBlock != 1 {
			t.Errorf("Expected skipBlock to be 1, but was %v for op %c %v", intr.skipBlock, op, intr.blockStack)
		}
	}
	intr.interpret(']')
	if intr.skipBlock != -1 {
		t.Errorf("Interpreter should be done skipping the current block.")
	}
}

func TestLoop(t *testing.T) {
	intr := NewInterpreter()
	loopCount := 3
	for i := 0; i < 3; i++ {
		intr.cells[i] = 2
	}
	intr.cells[3] = 1
	input := []byte{'[', '>', '-', ']'}
	for _, op := range input[1 : len(input)-1] {
		intr.interpret(op)
		if loopCount < 0 {
			t.Errorf("Looping too many times - expected only %v loops", loopCount)
		}
		loopCount--
	}

}

func TestInterpretHelloWorld(t *testing.T) {
	input := []byte("++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.")
	states := []struct {
		offset int
		index  int
		value  byte
		stack  []int
	}{
		{1, 0, 1, []int{}},       // +
		{2, 0, 2, []int{}},       // +
		{3, 0, 3, []int{}},       // +
		{4, 0, 4, []int{}},       // +
		{5, 0, 5, []int{}},       // +
		{6, 0, 6, []int{}},       // +
		{7, 0, 7, []int{}},       // +
		{8, 0, 8, []int{}},       // +
		{9, 0, 8, []int{8}},      // [
		{10, 1, 0, []int{8}},     // >
		{11, 1, 1, []int{8}},     // +
		{12, 1, 2, []int{8}},     // +
		{13, 1, 3, []int{8}},     // +
		{14, 1, 4, []int{8}},     // +
		{15, 1, 4, []int{8, 14}}, // [
		{16, 2, 0, []int{8, 14}}, // >
		{17, 2, 1, []int{8, 14}}, // +
		{18, 2, 2, []int{8, 14}}, // +
		{19, 3, 0, []int{8, 14}}, // >
		{20, 3, 1, []int{8, 14}}, // +
		{21, 3, 2, []int{8, 14}}, // +
		{22, 3, 3, []int{8, 14}}, // +
		{23, 4, 0, []int{8, 14}}, // >
		{24, 4, 1, []int{8, 14}}, // +
		{25, 4, 2, []int{8, 14}}, // +
		{26, 4, 3, []int{8, 14}}, // +
		{27, 5, 0, []int{8, 14}}, // >
		{28, 5, 1, []int{8, 14}}, // +
		{29, 4, 3, []int{8, 14}}, // <
		{30, 3, 3, []int{8, 14}}, // <
		{31, 2, 2, []int{8, 14}}, // <
		{32, 1, 4, []int{8, 14}}, // <
		// Good enough.
	}
	intr := NewInterpreter()
	errFmt := "Unexpected %v %v: '%+v'"
	for _, e := range states {
		intr.interpret(input[intr.offset])
		if e.offset != intr.offset {
			t.Errorf(errFmt, "offset", intr.offset, e)
		}
		if e.index != intr.index {
			t.Errorf(errFmt, "index", intr.index, e)
		}
		if e.value != intr.cells[intr.index] {
			t.Errorf(errFmt, "cell value", intr.cells[intr.index], e)
		}
		if len(e.stack) != len(intr.blockStack) {
			t.Errorf(errFmt, "block stack length", intr.blockStack, e)
		}
		for i, a := range intr.blockStack {
			if e.stack[i] != a {
				t.Errorf(errFmt, "block stack contents", intr.blockStack, e)
			}
		}
	}

}
