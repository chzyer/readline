package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

func usage(w io.Writer) {
	io.WriteString(w, `
sayhello: start to display oneline log per second
bye: quit
`[1:])
}

type Completer struct {
}

func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, off int) {
	list := [][]rune{
		[]rune("sayhello"), []rune("help"), []rune("bye"),
	}
	for i := 0; i <= 100; i++ {
		list = append(list, []rune(fmt.Sprintf("com%d", i)))
	}
	line = line[:pos]
	for _, r := range list {
		if strings.HasPrefix(string(r), string(line)) {
			newLine = append(newLine, r[len(line):])
		}
	}
	return newLine, len(line)
}

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:       "\033[31m»\033[0m ",
		HistoryFile:  "/tmp/readline.tmp",
		AutoComplete: new(Completer),
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err != nil {
			break
		}
		switch line {
		case "help":
			usage(l.Stderr())
		case "sayhello":
			go func() {
				for _ = range time.Tick(time.Second) {
					log.Println("hello")
				}
			}()
		case "bye":
			goto exit
		case "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
}
