package main

import (
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 readline.StaticPrompt("> "),
		HistoryFile:            "/tmp/readline-multiline",
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cmds = append(cmds, line)
		if !strings.HasSuffix(line, ";") {
			rl.SetPrompt(readline.StaticPrompt(">>> "))
			continue
		}
		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt(readline.StaticPrompt("> "))
		rl.SaveHistory(cmd)
		println(cmd)
	}
}
