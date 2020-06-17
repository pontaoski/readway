package app

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type DummyFile struct {
	Filename string
	f        *os.File
}

func (d DummyFile) Header() http.Header {
	return make(http.Header)
}

func (d DummyFile) WriteHeader(int) {

}

func (d *DummyFile) Write(b []byte) (int, error) {
	if d.f == nil {
		var err error
		d.f, err = os.Create(d.Filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	return d.f.Write(b)
}

type DummyByte struct {
	sb strings.Builder
}

func (d DummyByte) Header() http.Header {
	return make(http.Header)
}

func (d DummyByte) WriteHeader(int) {}

func (d *DummyByte) Write(b []byte) (int, error) {
	return d.sb.Write(b)
}

func (d DummyByte) String() string {
	return d.sb.String()
}
