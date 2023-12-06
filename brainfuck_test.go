package brainfuck

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	program, err := Compile("++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.")
	if err != nil {
		t.Error(err)
	}

	var output bytes.Buffer

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(&output)

	err = Execute(program, reader, writer)
	if err != nil {
		t.Error(err)
	}

	want := "Hello World!\n"
	if output.String() != want {
		t.Errorf("got: %s, want: %s", output.String(), want)
	}
}
