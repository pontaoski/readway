// +build wasm

package main

import _ "github.com/pontaoski/readway/app"

func main() {
	c := make(chan struct{}, 0)
	<-c
}
