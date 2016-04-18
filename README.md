# readline

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)
[![Build Status](https://travis-ci.org/chzyer/readline.svg?branch=master)](https://travis-ci.org/chzyer/readline)
[![GoDoc](https://godoc.org/github.com/chzyer/readline?status.svg)](https://godoc.org/github.com/chzyer/readline)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/chzyer/readline?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) 
[![OpenCollective](https://opencollective.com/readline/badge/backers.svg)](https://opencollective.com/readline#support)

Readline is A Pure Go Implementation of a libreadline-style Library.  
The goal is to be a powerful alternater for GNU-Readline.


**WHY:**
Readline will support most of features which GNU Readline is supported, and provide a pure go environment and a MIT license.

It can also provides shell-like interactives by using [flagly](https://github.com/chzyer/flagly) (demo: [flagly-shell](https://github.com/chzyer/flagly/blob/master/demo/flagly-shell/flagly-shell.go))

# Demo

![demo](https://raw.githubusercontent.com/chzyer/readline/assets/demo.gif)

Also works fine in windows

![demo windows](https://raw.githubusercontent.com/chzyer/readline/assets/windows.gif)


* [example/readline-demo](https://github.com/chzyer/readline/blob/master/example/readline-demo/readline-demo.go) The source code about the demo above

* [example/readline-im](https://github.com/chzyer/readline/blob/master/example/readline-im/readline-im.go) Example for how to write a IM program.

* [example/readline-multiline](https://github.com/chzyer/readline/blob/master/example/readline-multiline/readline-multiline.go) Example for how to parse command which can submit by multiple time.

* [example/readline-pass-strength](https://github.com/chzyer/readline/blob/master/example/readline-pass-strength/readline-pass-strength.go) A example about checking password strength, written by [@sahib](https://github.com/sahib)

# Todo
* Vi Mode is not completely finish
* More funny examples
* Support dumb/eterm-color terminal in emacs

# Features
* Support emacs/vi mode, almost all basic features that GNU-Readline is supported
* zsh-style backward/forward history search
* zsh-style completion
* Readline auto refresh when others write to Stdout while editing (it needs specify the Stdout/Stderr provided by *readline.Instance to others).
* Support colourful prompt in all platforms.

# Usage

* Import package

```
go get gopkg.in/readline.v1
```

or

```
go get github.com/chzyer/readline
```

* Simplest example

```go
import "gopkg.in/readline.v1"

rl, err := readline.New("> ")
if err != nil {
	panic(err)
}
defer rl.Close()

for {
	line, err := rl.Readline()
	if err != nil { // io.EOF, readline.ErrInterrupt
		break
	}
	println(line)
}
```

* Example with durable history

```go
rl, err := readline.NewEx(&readline.Config{
	Prompt: "> ",
	HistoryFile: "/tmp/readline.tmp",
})
if err != nil {
	panic(err)
}
defer rl.Close()

for {
	line, err := rl.Readline()
	if err != nil { // io.EOF, readline.ErrInterrupt
		break
	}
	println(line)
}
```

* Example with auto completion

```go
import (
	"gopkg.in/readline.v1"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("say",
		readline.PcItem("hello"),
		readline.PcItem("bye"),
	),
	readline.PcItem("help"),
)

rl, err := readline.NewEx(&readline.Config{
	Prompt:       "> ",
	AutoComplete: completer,
})
if err != nil {
	panic(err)
}
defer rl.Close()

for {
	line, err := rl.Readline()
	if err != nil { // io.EOF, readline.ErrInterrupt
		break
	}
	println(line)
}
```


# Shortcut

`Meta`+`B` means press `Esc` and `n` separately.  
Users can change that in terminal simulator(i.e. iTerm2) to `Alt`+`B`  
Notice: `Meta`+`B` is equals with `Alt`+`B` in windows.

* Shortcut in normal mode

| Shortcut           | Comment                                  |
|--------------------|------------------------------------------|
| `Ctrl`+`A`         | Beginning of line                        |
| `Ctrl`+`B` / `←`   | Backward one character                   |
| `Meta`+`B`         | Backward one word                        |
| `Ctrl`+`C`         | Send io.EOF                              |
| `Ctrl`+`D`         | Delete one character                     |
| `Meta`+`D`         | Delete one word                          |
| `Ctrl`+`E`         | End of line                              |
| `Ctrl`+`F` / `→`   | Forward one character                    |
| `Meta`+`F`         | Forward one word                         |
| `Ctrl`+`G`         | Cancel                                   |
| `Ctrl`+`H`         | Delete previous character                |
| `Ctrl`+`I` / `Tab` | Command line completion                  |
| `Ctrl`+`J`         | Line feed                                |
| `Ctrl`+`K`         | Cut text to the end of line              |
| `Ctrl`+`L`         | Clean screen (TODO)                      |
| `Ctrl`+`M`         | Same as Enter key                        |
| `Ctrl`+`N` / `↓`   | Next line (in history)                   |
| `Ctrl`+`P` / `↑`   | Prev line (in history)                   |
| `Ctrl`+`R`         | Search backwards in history              |
| `Ctrl`+`S`         | Search forwards in history               |
| `Ctrl`+`T`         | Transpose characters                     |
| `Meta`+`T`         | Transpose words (TODO)                   |
| `Ctrl`+`U`         | Cut text to the beginning of line        |
| `Ctrl`+`W`         | Cut previous word                        |
| `Backspace`        | Delete previous character                |
| `Meta`+`Backspace` | Cut previous word                        |
| `Enter`            | Line feed                                |


* Shortcut in Search Mode (`Ctrl`+`S` or `Ctrl`+`r` to enter this mode)

| Shortcut                | Comment                                     |
|-------------------------|---------------------------------------------|
| `Ctrl`+`S`              | Search forwards in history                  |
| `Ctrl`+`R`              | Search backwards in history                 |
| `Ctrl`+`C` / `Ctrl`+`G` | Exit Search Mode and revert the history     |
| `Backspace`             | Delete previous character                   |
| Other                   | Exit Search Mode                            |

* Shortcut in Complete Select Mode (double `Tab` to enter this mode)

| Shortcut                | Comment                                     |
|-------------------------|---------------------------------------------|
| `Ctrl`+`F`              | Move Forward                                |
| `Ctrl`+`B`              | Move Backward                               |
| `Ctrl`+`N`              | Move to next line                           |
| `Ctrl`+`P`              | Move to previous line                       |
| `Ctrl`+`A`              | Move to the first candicate in current line |
| `Ctrl`+`E`              | Move to the last candicate in current line  |
| `Tab` / `Enter`         | Use the word on cursor to complete          |
| `Ctrl`+`C` / `Ctrl`+`G` | Exit Complete Select Mode                   |
| Other                   | Exit Complete Select Mode                   |

# Tested with

| Environment                   | $TERM  |
|-------------------------------|--------|
| Mac OS X iTerm2               | xterm  |
| Mac OS X default Terminal.app | xterm  |
| Mac OS X iTerm2 Screen        | screen |
| Mac OS X iTerm2 Tmux          | screen |
| Ubuntu Server 14.04 LTS       | linux  |
| Centos 7                      | linux  |
| Windows 10                    | -      |

### Notice:
* `Ctrl`+`A` is not working in `screen` because it used as a control command by default

If you test it otherwhere, whether it works fine or not, please let me know!

## Who is using Readline

* [cockroachdb/cockroach](https://github.com/cockroachdb/cockroach)
* [youtube/doorman](https://github.com/youtube/doorman)
* [bom-d-van/harp](https://github.com/bom-d-van/harp)
* [abiosoft/ishell](https://github.com/abiosoft/ishell)
* [robertkrimen/otto](https://github.com/robertkrimen/otto)
* [Netflix/hal-9001](https://github.com/Netflix/hal-9001)

# Feedback

If you have any questions, please submit a github issue and any pull requests is welcomed :)

* [https://twitter.com/chzyer](https://twitter.com/chzyer)  
* [http://weibo.com/2145262190](http://weibo.com/2145262190)  


# Backers

Love Readline? Help me keep it alive by donating funds to cover project expenses!<br />
[[Become a backer](https://opencollective.com/readline)]

  <a href="https://opencollective.com/readline/backers/0/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/0/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/1/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/1/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/2/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/2/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/3/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/3/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/4/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/4/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/5/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/5/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/6/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/6/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/7/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/7/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/8/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/8/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/9/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/9/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/10/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/10/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/11/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/11/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/12/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/12/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/13/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/13/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/14/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/14/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/15/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/15/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/16/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/16/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/17/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/17/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/18/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/18/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/19/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/19/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/20/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/20/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/21/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/21/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/22/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/22/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/23/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/23/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/24/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/24/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/25/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/25/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/26/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/26/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/27/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/27/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/28/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/28/avatar">
  </a>
  <a href="https://opencollective.com/readline/backers/29/website" target="_blank">
    <img src="https://opencollective.com/readline/backers/29/avatar">
  </a>


