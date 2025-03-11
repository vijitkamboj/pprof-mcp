package profiler

import "fmt"

const InstructionsToInterpretResults = `
You are a expert profiler, result is structured by sampling type, for each sampling type, prepare a report, with insights about memory leaks, cpu usage, etc. give insights about memory leaks.
`

type ProfileType string

const (
	Heap         ProfileType = "heap"
	CPU          ProfileType = "profile"
	Goroutine    ProfileType = "goroutine"
	Threadcreate ProfileType = "threadcreate"
	Block        ProfileType = "block"
	Mutex        ProfileType = "mutex"
)

type ProfilerDepth int

const (
	LowLevelProfilerDepth ProfilerDepth = iota
	MediumLevelProfilerDepth
	HighLevelProfilerDepth
)

func (p ProfileType) String() string {
	return string(p)
}

func (p ProfileType) URL() string {
	return fmt.Sprintf("debug/pprof/%s", p)
}
