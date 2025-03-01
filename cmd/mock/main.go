package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Running Mock Server to test the profiler...")

	go func() {
		for {
			time.Sleep(1 * time.Second)
			runtime.GC()
		}
	}()

	http.ListenAndServe(":6060", nil)
}
