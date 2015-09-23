package readline

const (
	CharLineStart = 1
	CharBackward  = 2
	CharInterrupt = 3
	CharDelete    = 4
	CharLineEnd   = 5
	CharForward   = 6
	CharCannel    = 7
	CharCtrlH     = 8
	CharCtrlJ     = 10
	CharKill      = 11
	CharEnter     = 13
	CharNext      = 14
	CharPrev      = 16
	CharBckSearch = 18
	CharFwdSearch = 19
	CharTransform = 20
	CharCtrlW     = 23
	CharEsc       = 27
	CharEscapeEx  = 91
	CharBackspace = 127
)

const (
	MetaPrev = -iota - 1
	MetaNext
	MetaDelete
	MetaBackspace
	MetaTransform
)
