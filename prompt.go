package readline

type StaticPrompt string

func (s StaticPrompt) String() string {
	return string(s)
}

func NewStaticPrompt(s string) StaticPrompt {
	return StaticPrompt(s)
}
