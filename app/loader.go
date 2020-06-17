package app

import (
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
)

type Description struct {
	Summary string `xml:"summary,attr"`
	Body    string `xml:",chardata"`
}

type Anchorer struct {
	ID string
}

func (a *Anchorer) Anchor() string {
	if a.ID == "" {
		a.ID = strconv.FormatInt(rand.Int63(), 10)
	}
	return a.ID
}

type Argument struct {
	Anchorer
	Description Description `xml:"description"`
	Name        string      `xml:"name,attr"`
	Type        string      `xml:"type,attr"`
	Summary     string      `xml:"summary,attr"`
	Interface   string      `xml:"interface,attr"`
	AllowNull   bool        `xml:"allow-null,attr"`
	Enum        string      `xml:"enum,attr"`
}

func (a *Argument) Anchor() string {
	return a.Anchorer.Anchor()
}

type Request struct {
	Anchorer
	Description Description `xml:"description"`
	Arguments   []Argument  `xml:"arg"`
	Name        string      `xml:"name,attr"`
	Type        string      `xml:"type,attr"`
	Since       int         `xml:"since,attr"`
}

func (a *Request) Anchor() string {
	return a.Anchorer.Anchor()
}

type Event struct {
	Anchorer
	Description Description `xml:"description"`
	Arguments   []Argument  `xml:"arg"`
	Name        string      `xml:"name,attr"`
	Since       int         `xml:"since,attr"`
}

func (a *Event) Anchor() string {
	return a.Anchorer.Anchor()
}

type EnumEntry struct {
	Anchorer
	Description Description `xml:"description"`
	Name        string      `xml:"name,attr"`
	Summary     string      `xml:"summary,attr"`
	Value       string      `xml:"value,attr"`
	Since       int         `xml:"since,attr"`
}

func (a *EnumEntry) Anchor() string {
	return a.Anchorer.Anchor()
}

type Enum struct {
	Anchorer
	Description Description `xml:"description"`
	Entries     []EnumEntry `xml:"entry"`
	Name        string      `xml:"name,attr"`
	Since       int         `xml:"since,attr"`
	Bitfield    bool        `xml:"bitfield,attr"`
}

func (a *Enum) Anchor() string {
	return a.Anchorer.Anchor()
}

type Interface struct {
	Anchorer
	Name        string      `xml:"name,attr"`
	Version     int         `xml:"version,attr"`
	Description Description `xml:"description"`
	Requests    []Request   `xml:"request"`
	Events      []Event     `xml:"event"`
	Enums       []Enum      `xml:"enum"`
}

func (a *Interface) Anchor() string {
	return a.Anchorer.Anchor()
}

type Protocol struct {
	Anchorer
	Name        string      `xml:"name,attr"`
	Copyright   string      `xml:"copyright"`
	Description Description `xml:"description"`
	Interfaces  []Interface `xml:"interface"`
}

func (a *Protocol) Anchor() string {
	return a.Anchorer.Anchor()
}

func LoadProtocol(path string) Protocol {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	proto := Protocol{}
	err = xml.Unmarshal(data, &proto)
	if err != nil {
		panic(err)
	}
	return proto
}

func GetXMLs(path string) (ret []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xml") {
			ret = append(ret, filepath.Join(path, file.Name()))
		}
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
	stable := GetProtocols("protocols/stable")
	unstable := GetProtocols("protocols/unstable")
	plasma := GetProtocols("protocols/plasma")
	RenderProtocol(waylandMain)
	RenderProtocols(stable)
	RenderProtocols(unstable)
	RenderProtocols(plasma)
	RenderIndex(IndexData{waylandMain, stable, unstable, plasma})
}
