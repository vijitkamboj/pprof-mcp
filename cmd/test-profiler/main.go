package main

import (
	"fmt"

	"github.com/vijitkamboj/pprof-mcp/profiler"
)

func main() {
	pf := profiler.NewProfiler("localhost:6060", profiler.WithSamplers(
		profiler.NewHeapSampler(map[string]string{}),
		profiler.NewCPUSampler(map[string]string{
			"seconds": "5",
		}),
	))
	fmt.Println(pf.RunAll())

}
