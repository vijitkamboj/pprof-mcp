package profiler

import (
	"encoding/json"
	"log"
)

type Profiler struct {
	samplers []Sampler
	host     string
}

func NewProfiler(
	host string,
	opts ...func(*Profiler),
) *Profiler {
	p := &Profiler{
		host:     host,
		samplers: []Sampler{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithSamplers(samplers ...Sampler) func(*Profiler) {
	return func(p *Profiler) {
		p.samplers = samplers
	}
}

func (s *Profiler) RunAll() string {

	summaryPerSampler := map[ProfileType]interface{}{}

	for _, sp := range s.samplers {
		profile, err := GetParsedProfile(s.host, sp.Path(), sp.QueryParams())
		if err != nil {
			log.Fatalf("Failed to get profile: %v", err)
		}
		summary, err := sp.Summary(profile, LowLevelProfilerDepth)
		if err != nil {
			log.Fatalf("Failed to get summary: %v", err)
		}
		summaryPerSampler[sp.Name()] = summary
	}

	summaryPerSampler["InstructionsToInterpretResults"] = InstructionsToInterpretResults

	jsonData, err := json.Marshal(summaryPerSampler)
	if err != nil {
		log.Fatalf("Failed to marshal summary: %v", err)
	}

	return string(jsonData)

}
