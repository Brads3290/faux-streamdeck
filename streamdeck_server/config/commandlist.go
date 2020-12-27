package config

import (
	"encoding/xml"
)

// CommandBase list definition

type ButtonListSchema struct {
	XMLName xml.Name `xml:"buttons"`
	Buttons []Button `xml:"button"`
}

type Button struct {
	Commands []Command
	Name string `xml:"name,attr"`
	Icon string `xml:"icon,attr"`
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
				Cmd = new(ShortcutCommand)
			case "script":
				Cmd = new(ScriptCommand)
			case "command":
				Cmd = new(CommandCommand)
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



/*

 */
