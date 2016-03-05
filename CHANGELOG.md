# ChangeLog

### 1.2 - 2016-03-05

* Add a demo for checking password strength [example/readline-pass-strength](https://github.com/chzyer/readline/blob/master/example/readline-pass-strength/readline-pass-strength.go), , written by [@sahib](https://github.com/sahib)
* #23, support stdin remapping
* #27, add a `UniqueEditLine` to `Config`, which will erase the editing line after user submited it, usually use in IM.
* Add a demo for multiline [example/readline-multiline](https://github.com/chzyer/readline/blob/master/example/readline-multiline/readline-multiline.go) which can submit one SQL by multiple lines.
* Supports performs even stdin/stdout is not a tty.
* Add a new simple apis for single instance, check by [here](https://github.com/chzyer/readline/blob/master/std.go). It need to save history manually if using this api.
* #28, fixes the history is not working as expected.

### 1.1 - 2015-11-20

* #12 Add support for key `<Delete>`/`<Home>`/`<End>`
* Only enter raw mode as needed (calling `Readline()`), program will receive signal(e.g. Ctrl+C) if not interact with `readline`.
* Bugs fixed for `PrefixCompleter`
* Press `Ctrl+D` in empty line will cause `io.EOF` in error, Press `Ctrl+C` in anytime will cause `ErrInterrupt` instead of `io.EOF`, this will privodes a shell-like user experience.
* Customable Interrupt/EOF prompt in `Config`
* #17 Change atomic package to use 32bit function to let it runnable on arm 32bit devices
* Provides a new password user experience(`readline.ReadPasswordEx()`).

### 1.0 - 2015-10-14

* Initial public release.
