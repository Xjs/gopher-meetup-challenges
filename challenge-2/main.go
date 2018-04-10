package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type register string
type value int
type coprocessor map[register]value

func newCoprocessor(size int) coprocessor {
	c := make(map[register]value)
	const start = 'a'
	for i := 0; i < size; i++ {
		b := []byte{byte(start + i)}
		c[register(b)] = 0
	}
	return c
}

type input string

func (c coprocessor) value(i input) value {
	v, ok := c[register(i)]
	if ok {
		return v
	}

	v64, err := strconv.ParseInt(string(i), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Invalid value: %s", i))
	}
	return value(v64)
}

func (c coprocessor) set(x register, y input) coprocessor {
	c[x] = c.value(y)

	return c
}

func (c coprocessor) sub(x register, y input) coprocessor {
	c[x] -= c.value(y)

	return c
}

func (c coprocessor) mul(x register, y input) coprocessor {
	c[x] *= c.value(y)

	return c
}

func (c coprocessor) jnz(x input, y input) int {
	if c.value(x) == 0 {
		return 1
	}

	return int(c.value(y))
}

func main() {
	var instructions []string
	var size int
	var file string

	flag.IntVar(&size, "size", 8, "Size of our coprocessor (number of registers)")
	flag.StringVar(&file, "input", "-", "input file")
	flag.Parse()

	var raw []byte

	if file == "-" {
		var err error
		raw, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		raw, err = ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
	}

	rawInstructions := bytes.Split(raw, []byte("\n"))

	for _, instruction := range rawInstructions {
		s := strings.TrimSpace(string(instruction))
		if s != "" {
			instructions = append(instructions, s)
		}
	}

	var counter int

	c := newCoprocessor(size)

	var iter, mulCount int

	for {
		iter++
		if iter%100000 == 0 {
			fmt.Println(c)
		}
		if counter < 0 || counter >= len(instructions) {
			break
		}

		instruction := instructions[counter]

		crumbs := strings.Split(instruction, " ")
		name, pA, pB := crumbs[0], input(crumbs[1]), input(crumbs[2])

		switch name {
		case "set":
			reg := register(pA)
			inp := pB
			c = c.set(reg, inp)
		case "sub":
			reg := register(pA)
			inp := pB
			c = c.sub(reg, inp)
		case "mul":
			mulCount++
			reg := register(pA)
			inp := pB
			c = c.mul(reg, inp)
		case "jnz":
			counter += c.jnz(pA, pB)
			continue
		default:
			panic("Invalid instruction " + instruction)
		}

		counter++
	}

	fmt.Println(c)
	fmt.Println("mul: ", mulCount)
}
