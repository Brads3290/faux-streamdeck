package config


type Command interface {
	GetCommandType() string
}

type CommandBase struct {
	Type string `json:"type"`
}

func (cb CommandBase) GetCommandType() string {
	return cb.Type
}
