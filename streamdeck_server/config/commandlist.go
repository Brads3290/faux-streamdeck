package config

import (
	"encoding/xml"
)

// CommandBase list definition

type ButtonListSchema struct {
	XMLName xml.Name `xml:"buttons" json:"-"`
	Buttons []Button `xml:"button" json:"buttons"`
}

type Button struct {
	Commands []Command `json:"commands"`
	Id       string    `xml:"-" json:"id"`
	Name     string    `xml:"name,attr" json:"name"`
	Icon     string    `xml:"icon,attr" json:"icon"`
}

func (b *Button) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	b.Commands = make([]Command, 0)

	for _, v := range start.Attr {
		switch v.Name.Local {
		case "name":
			b.Name = v.Value
		case "icon":
			b.Icon = v.Value
		}
	}

	// decode inner elements
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		var Cmd Command = nil

		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "shortcut":
				Cmd = NewShortcutCommand()
			case "script":
				Cmd = NewScriptCommand()
			case "command":
				Cmd = NewCommandCommand()
			}

			// known child element found, decode it
			if Cmd != nil {
				err = d.DecodeElement(Cmd, &tt)
				if err != nil {
					return err
				}

				b.Commands = append(b.Commands, Cmd)
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}
}

type ShortcutCommand struct {
	CommandBase
	Keys string `xml:"keys,attr" json:"keys"`
}

func NewShortcutCommand() *ShortcutCommand {
	return &ShortcutCommand{
		CommandBase: CommandBase{Type: "shortcut"},
	}
}

type ScriptCommand struct {
	CommandBase
	Language string `xml:"language,attr" json:"language"`
	Path     string `xml:"path,attr" json:"path"`
}

func NewScriptCommand() *ScriptCommand {
	return &ScriptCommand{
		CommandBase: CommandBase{Type: "script"},
	}
}

type CommandCommand struct {
	CommandBase
	Text string `xml:"text,attr" json:"text"`
}

func NewCommandCommand() *CommandCommand {
	return &CommandCommand{
		CommandBase: CommandBase{Type: "command"},
	}
}
