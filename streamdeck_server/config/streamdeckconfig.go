package config

import "encoding/xml"

type StreamdeckConfigSchema struct {
	XMLName xml.Name `xml:"streamdeck"`
	Server StreamdeckServer `xml:"server"`
}

type StreamdeckServer struct {
	ListenOn string `xml:"listen,attr"`
}
