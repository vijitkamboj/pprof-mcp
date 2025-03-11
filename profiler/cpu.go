package profiler

import "github.com/google/pprof/profile"

type CPUSampler struct {
	queryParams map[string]string
}

func NewCPUSampler(queryParams map[string]string) *CPUSampler {
	return &CPUSampler{queryParams: queryParams}
}

func (s *CPUSampler) Path() string {
	return CPU.URL()
}

func (s *CPUSampler) Name() ProfileType {
	return CPU
}

func (s *CPUSampler) QueryParams() map[string]string {
	return s.queryParams
}

func (s *CPUSampler) Summary(prof *profile.Profile, depth ProfilerDepth) (any, error) {
	return prof, nil
}
