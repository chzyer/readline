package main

import (
	"fmt"
	"strings"

	"github.com/gohxs/readline"
)

type TestCompleter struct {
}

func (t *TestCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	res := [][]rune{}
	for i := 0; i < 2000; i++ {
		sample := []rune(fmt.Sprintf("%04d", i+1000))
		if strings.HasPrefix(string(sample), string(line)) {
			res = append(res, sample[len(line):])
		}
	}
	return res, pos

	// Ingore for now
	// Limiter merger whatever?
	in := len(res[0]) - 1
	for len(res) > 10 && in > 0 { // Match only the starts
		newRes := [][]rune{}
		start := res[0][:in] // add to start
		newRes = append(newRes, start)
		for _, v := range res[1:] {
			if contains(newRes, v[:in]) {
				continue
			}
			newRes = append(newRes, v[:in])
		}
		in--
		res = newRes
	}
	return res, pos
}

func contains(list [][]rune, search []rune) bool {
	for _, v := range list {
		if string(v) == string(search) {
			return true
		}
	}
	return false
}

func main() {

	rl, _ := readline.NewEx(&readline.Config{
		Prompt:           "test> ",
		AutoComplete:     &TestCompleter{},
		MaxCompleteLines: 3,
	})

	for {
		_, err := rl.Readline()
		if err != nil {
			return
		}

	}

}
