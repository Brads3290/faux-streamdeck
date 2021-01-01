package config

type Command interface {
	GetCommandType() string
}

type CommandBase struct {
	Type string `json:"type"`
}

// GetCommandType returns the type of command that is represented by this object
func (cb CommandBase) GetCommandType() string {
	return cb.Type
}
