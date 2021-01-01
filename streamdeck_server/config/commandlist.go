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
	Commands []Command `json:"-"`
	Id       string    `xml:"-" json:"id"`
	Name     string    `xml:"name,attr" json:"name"`
	Icon     string    `xml:"icon,attr" json:"icon"`
}

// UnmarshalXML is used for Button to create the command list, because the commands can be one of several
// different types, each needing different representions.
func (b *Button) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	b.Commands = make([]Command, 0)

	// Decode the attributes for the Button object itself
	for _, v := range start.Attr {
		switch v.Name.Local {
		case "name":
			b.Name = v.Value
		case "icon":
			b.Icon = v.Value
		}
	}

	// Decode the inner elements, which could be one of multiple types:
	// - shortcut
	// - script
	// - shell command
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		var Cmd Command = nil

		// For each token, check if it's a start element or an end element, otherwise, ignore.
		switch tt := t.(type) {
		case xml.StartElement:

			// If start element, decode it based on it's element name
			switch tt.Name.Local {
			case "shortcut":
				Cmd = NewShortcutCommand()
			case "script":
				Cmd = NewScriptCommand()
			case "command":
				Cmd = NewShellCommand()
			default:
				Cmd = nil
			}

			// Known child element found, decode it
			if Cmd != nil {
				err = d.DecodeElement(Cmd, &tt)
				if err != nil {
					return err
				}

				b.Commands = append(b.Commands, Cmd)
			}
		case xml.EndElement:

			// If this end element is the end element for the <button> tag that we're
			// processing, we're done!
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

type ShellCommand struct {
	CommandBase
	Text string `xml:"text,attr" json:"text"`
}

func NewShellCommand() *ShellCommand {
	return &ShellCommand{
		CommandBase: CommandBase{Type: "command"},
	}
}
