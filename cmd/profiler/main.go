package main

import (
	"encoding/json"
	"log"

	"context"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/vijitkamboj/pprof-mcp/sampler"
)

func main() {
	done := make(chan struct{})

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
	err := server.RegisterTool("profiler", "Get the performance profile of a golang server [ always take host as input from user first, if user denies, do not proceed with this tool ]", profileHandler)
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	<-done
}

// >>>>>>>>>>>>>>>>>>>>>>>>> PROFILER <<<<<<<<<<<<<<<<<<<<<<<<<<

type Profiler struct {
	samplers []sampler.Sampler
	host     string
}

func NewProfiler(
	host string,
) *Profiler {
	return &Profiler{
		host:     host,
		samplers: []sampler.Sampler{},
	}
}

func (s *Profiler) addSamplers(samplers ...sampler.Sampler) {
	s.samplers = append(s.samplers, samplers...)
}

func (s *Profiler) profile() string {

	summaryPerSampler := map[string]interface{}{}
	for _, sp := range s.samplers {
		profile, err := sampler.GetParsedProfile(s.host, sp.Path(), sp.QueryParams())
		if err != nil {
			log.Fatalf("Failed to get profile: %v", err)
		}
		summary, err := sp.Summary(profile)
		if err != nil {
			log.Fatalf("Failed to get summary: %v", err)
		}
		summaryPerSampler[sp.Path()] = summary

	}

	summaryPerSampler["INSTRUCTIONS_TO_INTERPRET_THE_RESULTS"] = "You are a expert profiler, result is structured by sampling type, for each sampling type, prepare a report, with insights about memory leaks, cpu usage, etc. give insights about memory leaks."

	jsonData, err := json.Marshal(summaryPerSampler)
	if err != nil {
		log.Fatalf("Failed to marshal summary: %v", err)
	}

	return string(jsonData)

}

type ProfileRequest struct {
	Host string `json:"host" jsonschema:"required,description="`
}

func profileHandler(ctx context.Context, request ProfileRequest) (*mcp_golang.ToolResponse, error) {

	profiler := NewProfiler(request.Host)
	profiler.addSamplers(
		sampler.NewHeapSampler(map[string]string{}),
	)
	data := profiler.profile()
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(data)), nil

}
