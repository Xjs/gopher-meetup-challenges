package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type programs []byte

func (p programs) fit(n int) int {
	n = n % len(p)
	if n < 0 {
		n = n + len(p)
	}

	return n
}

func (p programs) String() string {
	return string(p)
}

func (p programs) spin(n int) programs {
	index := p.fit(len(p) - n)
	return append(p[index:], p[:index]...)
}

func (p programs) exchange(a, b int) programs {
	a = p.fit(a)
	b = p.fit(b)

	p[b], p[a] = p[a], p[b]

	return p
}

func (p programs) partner(a, b byte) programs {
	var indexA, indexB int

	indexA = -1
	indexB = -1

	for i, prog := range p {
		if prog == a {
			indexA = i
		}

		if prog == b {
			indexB = i
		}
	}

	if indexA < 0 || indexB < 0 {
		panic(fmt.Sprint("Programs ", a, " and ", b, " not found in ", p))
	}

	return p.exchange(indexA, indexB)
}

func mustParseInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("%s is no integer: %v", s, err))
	}
	return int(i)
}

func (p programs) parseCommand(command string) programs {
	switch command[0] {
	case 's':
		return p.spin(mustParseInt(command[1:]))
	case 'x':
		ns := strings.Split(command[1:], "/")
		if len(ns) != 2 {
			panic("Invalid number of arguments")
		}
		a := mustParseInt(ns[0])
		b := mustParseInt(ns[1])

		return p.exchange(a, b)
	case 'p':
		ns := strings.Split(command[1:], "/")
		if len(ns) != 2 {
			panic("Invalid number of arguments")
		}

		a := ns[0][0]
		b := ns[1][0]

		return p.partner(a, b)
	default:
		panic(fmt.Sprint("Unknown command: ", command))
	}
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Output state after each step")

	commands, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	p := programs("abcdefghijklmnop")

	if debug {
		fmt.Println("start - ", p)
	}

	for _, command := range strings.Split(string(commands), ",") {
		p = p.parseCommand(command)
		if debug {
			fmt.Println(command, " - ", p)
		}
	}

	fmt.Println(p)
}
