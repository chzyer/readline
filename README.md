[![Build Status](https://travis-ci.org/chzyer/readline.svg?branch=master)](https://travis-ci.org/chzyer/readline)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)
[![Version](https://img.shields.io/github/tag/chzyer/readline.svg)](https://github.com/chzyer/readline/releases)
[![GoDoc](https://godoc.org/github.com/chzyer/readline?status.svg)](https://godoc.org/github.com/chzyer/readline)
[![OpenCollective](https://opencollective.com/readline/badge/backers.svg)](#backers)
[![OpenCollective](https://opencollective.com/readline/badge/sponsors.svg)](#sponsors)

<p align="center">
<img src="http://0xdf.com:39178/logo.png?1" />
<a href="https://asciinema.org/a/32oseof9mkilg7t7d4780qt4m" target="_blank"><img src="https://asciinema.org/a/32oseof9mkilg7t7d4780qt4m.png" width="654"/></a>
<img src="http://0xdf.com:39178/logo_foot.png?1" />
</p>


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



# Tested with

| Environment                   | $TERM  |
| ----------------------------- | ------ |
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
* [remind101/empire](https://github.com/remind101/empire)
* [youtube/doorman](https://github.com/youtube/doorman)
* [bom-d-van/harp](https://github.com/bom-d-van/harp)
* [abiosoft/ishell](https://github.com/abiosoft/ishell)
* [robertkrimen/otto](https://github.com/robertkrimen/otto)
* [Netflix/hal-9001](https://github.com/Netflix/hal-9001)
* [docker/go-p9p](https://github.com/docker/go-p9p)
* [mehrdadrad/mylg](https://github.com/mehrdadrad/mylg)

# Feedback

If you have any questions, please submit a github issue and any pull requests is welcomed :)

* [https://twitter.com/chzyer](https://twitter.com/chzyer)  
* [http://weibo.com/2145262190](http://weibo.com/2145262190)  


# Backers

Love Readline? Help me keep it alive by donating funds to cover project expenses!<br />
[[Become a backer](https://opencollective.com/readline#backer)]

<a href="https://opencollective.com/readline/backer/0/website" target="_blank"><img src="https://opencollective.com/readline/backer/0/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/1/website" target="_blank"><img src="https://opencollective.com/readline/backer/1/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/2/website" target="_blank"><img src="https://opencollective.com/readline/backer/2/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/3/website" target="_blank"><img src="https://opencollective.com/readline/backer/3/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/4/website" target="_blank"><img src="https://opencollective.com/readline/backer/4/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/5/website" target="_blank"><img src="https://opencollective.com/readline/backer/5/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/6/website" target="_blank"><img src="https://opencollective.com/readline/backer/6/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/7/website" target="_blank"><img src="https://opencollective.com/readline/backer/7/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/8/website" target="_blank"><img src="https://opencollective.com/readline/backer/8/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/9/website" target="_blank"><img src="https://opencollective.com/readline/backer/9/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/10/website" target="_blank"><img src="https://opencollective.com/readline/backer/10/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/11/website" target="_blank"><img src="https://opencollective.com/readline/backer/11/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/12/website" target="_blank"><img src="https://opencollective.com/readline/backer/12/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/13/website" target="_blank"><img src="https://opencollective.com/readline/backer/13/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/14/website" target="_blank"><img src="https://opencollective.com/readline/backer/14/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/15/website" target="_blank"><img src="https://opencollective.com/readline/backer/15/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/16/website" target="_blank"><img src="https://opencollective.com/readline/backer/16/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/17/website" target="_blank"><img src="https://opencollective.com/readline/backer/17/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/18/website" target="_blank"><img src="https://opencollective.com/readline/backer/18/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/19/website" target="_blank"><img src="https://opencollective.com/readline/backer/19/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/20/website" target="_blank"><img src="https://opencollective.com/readline/backer/20/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/21/website" target="_blank"><img src="https://opencollective.com/readline/backer/21/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/22/website" target="_blank"><img src="https://opencollective.com/readline/backer/22/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/23/website" target="_blank"><img src="https://opencollective.com/readline/backer/23/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/24/website" target="_blank"><img src="https://opencollective.com/readline/backer/24/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/25/website" target="_blank"><img src="https://opencollective.com/readline/backer/25/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/26/website" target="_blank"><img src="https://opencollective.com/readline/backer/26/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/27/website" target="_blank"><img src="https://opencollective.com/readline/backer/27/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/28/website" target="_blank"><img src="https://opencollective.com/readline/backer/28/avatar.svg"></a>
<a href="https://opencollective.com/readline/backer/29/website" target="_blank"><img src="https://opencollective.com/readline/backer/29/avatar.svg"></a>


# Sponsors

Become a sponsor and get your logo here on our Github page. [[Become a sponsor](https://opencollective.com/readline#sponsor)]

<a href="https://opencollective.com/readline/sponsor/0/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/1/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/2/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/3/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/4/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/5/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/6/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/7/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/8/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/9/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/9/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/10/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/10/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/11/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/11/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/12/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/12/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/13/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/13/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/14/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/14/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/15/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/15/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/16/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/16/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/17/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/17/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/18/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/18/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/19/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/19/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/20/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/20/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/21/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/21/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/22/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/22/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/23/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/23/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/24/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/24/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/25/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/25/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/26/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/26/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/27/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/27/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/28/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/28/avatar.svg"></a>
<a href="https://opencollective.com/readline/sponsor/29/website" target="_blank"><img src="https://opencollective.com/readline/sponsor/29/avatar.svg"></a>

