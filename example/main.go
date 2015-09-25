package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/traetox/readline"
)

func usage(w io.Writer) {
	io.WriteString(w, `
sayhello: start to display oneline log per second
bye: quit
`[1:])
}

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:      "\033[31m»\033[0m ",
		HistoryFile: "/tmp/readline.tmp",
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
		case "write":
			fmt.Fprintf(l, "look, i am writer too\n")
		case "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
}
