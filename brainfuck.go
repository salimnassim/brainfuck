package brainfuck

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type instruction struct {
	operator uint16
	operand  uint16
}

const (
	opPtrRight = iota
	opPtrLeft
	opValInc
	opValDec
	opOut
	opIn
	opJmpFw
	opJmpBk
)

type Brainfucker interface {
	Compile(input string) (program []instruction, err error)
	Execute(program []instruction, reader io.Reader, writer io.Writer) error
}

func Compile(input string) (program []instruction, err error) {
	var cnt, jmpCnt uint16
	jmpStack := make([]uint16, 0)

	for _, c := range input {
		switch c {
		case '>':
			program = append(program, instruction{opPtrRight, 0})
		case '<':
			program = append(program, instruction{opPtrLeft, 0})
		case '+':
			program = append(program, instruction{opValInc, 0})
		case '-':
			program = append(program, instruction{opValDec, 0})
		case '.':
			program = append(program, instruction{opOut, 0})
		case ',':
			program = append(program, instruction{opIn, 0})
		case '[':
			program = append(program, instruction{opJmpFw, 0})
			jmpStack = append(jmpStack, cnt)
		case ']':
			if len(jmpStack) == 0 {
				return nil, errors.New("unmatched ']'")
			}
			jmpCnt = jmpStack[len(jmpStack)-1]
			jmpStack = jmpStack[:len(jmpStack)-1]
			program = append(program, instruction{opJmpBk, jmpCnt})
			program[jmpCnt].operand = cnt
		default:
			cnt--
		}
		cnt++
	}

	if len(jmpStack) != 0 {
		return nil, errors.New("unmatched '['")
	}

	return program, nil
}

func Execute(program []instruction, reader io.Reader, writer io.Writer) error {
	mem := make([]uint16, 65536)
	var memPtr uint16 = 0

	bufReader := bufio.NewReader(reader)
	bufWriter := bufio.NewWriter(writer)

	for pc := 0; pc < len(program); pc++ {
		switch program[pc].operator {
		case opPtrRight:
			memPtr++
		case opPtrLeft:
			memPtr--
		case opValInc:
			mem[memPtr]++
		case opValDec:
			mem[memPtr]--
		case opOut:
			bufWriter.Write([]byte(fmt.Sprintf("%c", mem[memPtr])))
		case opIn:
			val, err := bufReader.ReadByte()
			if err != nil {
				return err
			}
			mem[memPtr] = uint16(val)
		case opJmpFw:
			if mem[memPtr] == 0 {
				pc = int(program[pc].operand)
			}
		case opJmpBk:
			if mem[memPtr] > 0 {
				pc = int(program[pc].operand)
			}
		default:
			return errors.New("unknown operator")
		}
	}

	bufWriter.Flush()

	return nil
}
