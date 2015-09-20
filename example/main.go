package main

import (
	"io"
	"log"
	"strconv"
	"time"

	"github.com/chzyer/readline"
)

func main() {
	t, err := readline.NewTerminal()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	l := t.NewReadline("> ")
	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err != nil {
			break
		}
		switch line {
		case "help":
			io.WriteString(l.Stderr(), "sayhello: start to display oneline log per second\nbye: quit\n")
		case "sayhello":
			go func() {
				for _ = range time.Tick(time.Second) {
					log.Println("hello")
				}
			}()
		case "bye":
			goto exit
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
exit:
}
