package config


type Command interface {
	GetCommandType() string
}

type CommandBase struct {
	Type string
}

func (cb CommandBase) GetCommandType() string {
	return cb.Type
}

type ShortcutCommand struct {
	CommandBase
	Keys string `xml:"keys,attr"`

}

type ScriptCommand struct {
	CommandBase
	Language string `xml:"language,attr"`
	Path string `xml:"path,attr"`
}

type CommandCommand struct {
	CommandBase
	Text string `xml:"text,attr"`
}