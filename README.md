# readline
A pure go implementation for gnu readline.

![demo](https://raw.githubusercontent.com/chzyer/readline/master/example/demo.gif)

# usage

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



