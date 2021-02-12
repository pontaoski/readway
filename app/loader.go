package app

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var rplcr = strings.NewReplacer(
	" ", "_",
)

func toAnchor(s string) string {
	return rplcr.Replace(strings.ToLower(s))
}

type Description struct {
	Summary string `xml:"summary,attr"`
	Body    string `xml:",chardata"`
}

type Anchorer struct {
	ID string
}

type Argument struct {
	Description Description `xml:"description"`
	Name        string      `xml:"name,attr"`
	Type        string      `xml:"type,attr"`
	Summary     string      `xml:"summary,attr"`
	Interface   string      `xml:"interface,attr"`
	AllowNull   bool        `xml:"allow-null,attr"`
	Enum        string      `xml:"enum,attr"`
}

type Request struct {
	Parent      *Interface
	Description Description `xml:"description"`
	Arguments   []Argument  `xml:"arg"`
	Name        string      `xml:"name,attr"`
	Type        string      `xml:"type,attr"`
	Since       int         `xml:"since,attr"`
}

func (a *Request) Anchor() string {
	return a.Parent.Anchor() + "_" + toAnchor(a.Name)
}

type Event struct {
	Parent      *Interface
	Description Description `xml:"description"`
	Arguments   []Argument  `xml:"arg"`
	Name        string      `xml:"name,attr"`
	Since       int         `xml:"since,attr"`
}

func (a *Event) Anchor() string {
	return a.Parent.Anchor() + "_" + toAnchor(a.Name)
}

type EnumEntry struct {
	Parent      *Enum
	Description Description `xml:"description"`
	Name        string      `xml:"name,attr"`
	Summary     string      `xml:"summary,attr"`
	Value       string      `xml:"value,attr"`
	Since       int         `xml:"since,attr"`
}

func (a *EnumEntry) Anchor() string {
	return a.Parent.Anchor() + "_" + toAnchor(a.Name)
}

type Enum struct {
	Parent      *Interface
	Description Description  `xml:"description"`
	Entries     []*EnumEntry `xml:"entry"`
	Name        string       `xml:"name,attr"`
	Since       int          `xml:"since,attr"`
	Bitfield    bool         `xml:"bitfield,attr"`
}

func (a *Enum) Anchor() string {
	return a.Parent.Anchor() + "_" + toAnchor(a.Name)
}

type Interface struct {
	Name        string      `xml:"name,attr"`
	Version     int         `xml:"version,attr"`
	Description Description `xml:"description"`
	Requests    []*Request  `xml:"request"`
	Events      []*Event    `xml:"event"`
	Enums       []*Enum     `xml:"enum"`
}

func (a *Interface) Anchor() string {
	return toAnchor(a.Name)
}

type Protocol struct {
	Name        string       `xml:"name,attr"`
	Copyright   string       `xml:"copyright"`
	Description Description  `xml:"description"`
	Interfaces  []*Interface `xml:"interface"`
}

func (p *Protocol) SetTree() {
	for _, iface := range p.Interfaces {
		for _, enum := range iface.Enums {
			enum.Parent = iface
			for _, entry := range enum.Entries {
				entry.Parent = enum
			}
		}
		for _, ev := range iface.Events {
			ev.Parent = iface
		}
		for _, req := range iface.Requests {
			req.Parent = iface
		}
	}
}

func LoadProtocol(path string) Protocol {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	proto := Protocol{}
	err = xml.Unmarshal(data, &proto)
	proto.SetTree()
	if err != nil {
		panic(err)
	}
	return proto
}

func GetXMLs(path string) (ret []string) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".xml") {
			ret = append(ret, path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func GetProtocols(path string) (ret []Protocol) {
	files := GetXMLs(path)
	for _, file := range files {
		ret = append(ret, LoadProtocol(file))
	}
	return
}

func RenderProtocols(protocols []Protocol) {
	for _, proto := range protocols {
		RenderProtocol(proto)
	}
}

type IndexData struct {
	Wayland  Protocol
	Stable   []Protocol
	Unstable []Protocol
	Plasma   []Protocol
}

func LoadProtocols() {
	waylandMain := LoadProtocol("protocols/wayland.xml")
	stable := GetProtocols("protocols/stable/")
	unstable := GetProtocols("protocols/unstable/")
	plasma := GetProtocols("protocols/plasma/")
	RenderProtocol(waylandMain)
	RenderProtocols(stable)
	RenderProtocols(unstable)
	RenderProtocols(plasma)
	RenderIndex(IndexData{waylandMain, stable, unstable, plasma})
}
