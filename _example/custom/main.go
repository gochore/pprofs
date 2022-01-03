package main

import (
	"log"
	"time"

	"github.com/gochore/pprofs"
)

func main() {
	if err := pprofs.EnableCapture(
		pprofs.WithProfiles(
			pprofs.CpuProfile().WithDuration(10*time.Second),
			pprofs.HeapProfile(),
			pprofs.MutexProfile(),
			pprofs.BlockProfile().WithRate(1),
			pprofs.GoroutineProfile(),
			pprofs.ThreadcreateProfile(),
		),
		pprofs.WithTrigger(pprofs.NewIntervalTrigger(15*time.Second)),
		pprofs.WithStorage(pprofs.NewFileStorage("custom", "/tmp/pprofs/", time.Hour)),
		pprofs.WithLogger(log.Default()),
	); err != nil {
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
