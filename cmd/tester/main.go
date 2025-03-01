package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vijitkamboj/pprof-mcp/sampler"
)

func main() {
	profiler := NewProfiler("localhost:6060")
	profiler.addSamplers(
		sampler.NewHeapSampler(map[string]string{}),
	)
	data := profiler.profile()

	fmt.Println(data)

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

	jsonData, err := json.Marshal(summaryPerSampler)
	if err != nil {
		log.Fatalf("Failed to marshal summary: %v", err)
	}

	return string(jsonData)

}
