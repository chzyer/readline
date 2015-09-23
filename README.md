# readline

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)
[![Build Status](https://travis-ci.org/chzyer/readline.svg?branch=master)](https://travis-ci.org/chzyer/readline)
[![GoDoc](https://godoc.org/github.com/chzyer/readline?status.svg)](https://godoc.org/github.com/chzyer/readline)  

A pure go implementation for gnu readline.

# Demo

![demo](https://raw.githubusercontent.com/chzyer/readline/master/example/demo.gif)

# Shortcut

> `Meta + n` means press `Esc` and `n` separately

| Shortcut           | Comment                           | Support |
|--------------------|-----------------------------------|---------|
| `Ctrl`+`A`         | Beginning of line                 | Yes     |
| `Ctrl`+`B` / `←`   | Backward one character            | Yes     |
| `Meta`+`B`         | Backward one word                 | Yes     |
| `Ctrl`+`C`         | Send io.EOF                       | Yes     |
| `Ctrl`+`D`         | Delete one character              | Yes     |
| `Meta`+`D`         | Delete one word                   | Yes     |
| `Ctrl`+`E`         | End of line                       | Yes     |
| `Ctrl`+`F` / `→`   | Forward one character             | Yes     |
| `Meta`+`F`         | Forward one word                  | Yes     |
| `Ctrl`+`G`         | Cancel                            | Yes     |
| `Ctrl`+`H`         | Delete previous character         | Yes     |
| `Ctrl`+`I` / `Tab` | Command line completion           | No      |
| `Ctrl`+`J`         | Line feed                         | Yes     |
| `Ctrl`+`K`         | Cut text to the end of line       | Yes     |
| `Ctrl`+`L`         | Clean screen                      | No      |
| `Ctrl`+`M`         | Same as Enter key                 | Yes     |
| `Ctrl`+`N` / `↓`   | Next line (in history)            | Yes     |
| `Ctrl`+`P` / `↑`   | Prev line (in history)            | Yes     |
| `Ctrl`+`Q`         | Resume transmission               | No      |
| `Ctrl`+`R`         | Search backwards in history       | Yes     |
| `Ctrl`+`S`         | Search forwards in history        | Yes     |
| `Ctrl`+`T`         | Transpose characters              | Yes     |
| `Ctrl`+`U`         | Cut text to the beginning of line | No      |
| `Ctrl`+`W`         | Cut previous word                 | Yes     |
| `Backspace`        | Delete previous character         | Yes     |
| `Meta`+`Backspace` | Cut previous word                 | Yes     |
| `Enter`            | Line feed                         | Yes     |

