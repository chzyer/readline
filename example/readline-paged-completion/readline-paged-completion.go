package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// A completor that will give a lot of completions for showcasing the paging functionality
type Completor struct{}

func (c *Completor) Do(line []rune, pos int) ([][]rune, int) {
	completion := make([][]rune, 0, 10000)
	for i := range 10000 {
		s := fmt.Sprintf("%s%020d", randSeq(1), i)
		completion = append(completion, []rune(s))
	}
	return completion, pos
}

func main() {
	c := Completor{}
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		AutoComplete:    &c,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
}
