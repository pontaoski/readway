package app

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/qor/render"
)

var funcMap = template.FuncMap{
	"unescape": func(in string) template.HTML {
		return template.HTML(in)
	},
	"isLast": func(index, len int) bool {
		return index+1 == len
	},
	"hasField": func(v interface{}, name string) bool {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() != reflect.Struct {
			return false
		}
		return rv.FieldByName(name).IsValid()
	},
	"join_strings": strings.Join,
}

var renderer = render.New(&render.Config{
	ViewPaths:     []string{"app/templates"},
	DefaultLayout: "application",
	FuncMapMaker: func(*render.Render, *http.Request, http.ResponseWriter) template.FuncMap {
		return funcMap
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
