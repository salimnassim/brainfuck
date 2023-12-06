package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/salimnassim/brainfuck"
)

func main() {
	var inputFile = flag.String("file", "", "path to input file")
	flag.Parse()

	if *inputFile == "" {
		fmt.Print("no input file specified")
		return
	}

	f, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("unable to open input file: %s", err)
		return
	}

	inputBuffer := &bytes.Buffer{}
	inputReader := bufio.NewReader(f)
	inputBuffer.ReadFrom(inputReader)

	if len(inputBuffer.String()) == 0 {
		fmt.Print("input file is empty")
		return
	}

	p, err := brainfuck.Compile(inputBuffer.String())
	if err != nil {
		fmt.Printf("compile error: %s", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	err = brainfuck.Execute(p, reader, writer)
	if err != nil {
		fmt.Printf("execute error: %s", err)
		return
	}
}
