package main

import (
	"time"

	"github.com/gochore/pprofs"
)

func main() {
	if err := pprofs.EnableCapture(); err != nil {
		panic(err)
	}
	run()
}

func run() {
	for {
		time.Sleep(time.Second)
		for i := 0; i < 100_0000; i++ {
		}
		_ = make([]byte, 16*1024*1024)
	}
}
