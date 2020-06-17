package app

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/qor/render"
)

var renderer = render.New(&render.Config{
	ViewPaths:     []string{"app/templates"},
	DefaultLayout: "application",
	FuncMapMaker: func(*render.Render, *http.Request, http.ResponseWriter) template.FuncMap {
		return map[string]interface{}{
			"unescape": func(in string) template.HTML {
				return template.HTML(in)
			},
			"join_strings": strings.Join,
		}
	},
})

func RenderProtocol(protocol Protocol) {
	println("Rendering " + protocol.Name + "...")
	err := renderer.Execute("protocol", protocol, &http.Request{}, &DummyFile{"out/" + protocol.Name + ".html", nil})
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func RenderProtocolToString(protocol Protocol) string {
	dummy := DummyByte{}
	err := renderer.Execute("protocol", protocol, &http.Request{}, &dummy)
	if err != nil {
		return err.Error()
	}
	return dummy.String()
}

func RenderIndex(data IndexData) {
	println("Rendering Index...")
	err := renderer.Execute("index", data, &http.Request{}, &DummyFile{"out/index.html", nil})
	if err != nil {
		log.Println(err.Error())
		return
	}
}
