package main

import (
	"io"
	"log"
	"strconv"
	"time"

	"github.com/chzyer/readline"
)

func usage(w io.Writer) {
	io.WriteString(w, `
sayhello: start to display oneline log per second
bye: quit
`[1:])
}

func main() {
	l, err := readline.New("home -> ")
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
