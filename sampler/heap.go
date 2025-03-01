package sampler

import (
	"sort"

	"github.com/google/pprof/profile"
)

type HeapSampler struct {
	queryParams map[string]string
}

func NewHeapSampler(queryParams map[string]string) *HeapSampler {
	return &HeapSampler{
		queryParams: queryParams,
	}
}

func (s *HeapSampler) Summary(prof *profile.Profile) (any, error) {
	return createSummary(prof)
}

func (s *HeapSampler) Path() string {
	return "/debug/pprof/heap"
}

func (s *HeapSampler) QueryParams() map[string]string {
	return s.queryParams
}

// Summary represents a condensed profile summary suitable for LLM analysis
type HeapSummary struct {
	SampleTypes []SampleTypeInfo     `json:"sampleTypes"`
	TopSamples  []CondensedSample    `json:"topSamples"`
	Stats       map[string]int64     `json:"stats"`
	MemoryUsage []MemoryUsageSummary `json:"memoryUsage"`
}

type SampleTypeInfo struct {
	Type string `json:"type"`
	Unit string `json:"unit"`
}

type CondensedSample struct {
	Value        []int64   `json:"value"`
	TopFunctions []string  `json:"topFunctions"`
	Labels       LabelInfo `json:"labels,omitempty"`
}

type LabelInfo struct {
	Labels   map[string]string  `json:"labels,omitempty"`
	NumLabel map[string][]int64 `json:"numLabels,omitempty"`
}

type MemoryUsageSummary struct {
	Size  int64  `json:"size"`
	Count int64  `json:"count"`
	Stack string `json:"topFunction"`
}

func createSummary(prof *profile.Profile) (any, error) {
	summary := HeapSummary{
		Stats:       make(map[string]int64),
		MemoryUsage: []MemoryUsageSummary{},
	}

	// Extract sample types
	for _, sampleType := range prof.SampleType {
		summary.SampleTypes = append(summary.SampleTypes, SampleTypeInfo{
			Type: sampleType.Type,
			Unit: sampleType.Unit,
		})
	}

	// Find the index for inuse_space or alloc_space
	spaceIndex := findMemorySpaceIndex(prof)
	objectsIndex := findMemoryObjectsIndex(prof)

	// Summary stats
	summary.Stats["totalSamples"] = int64(len(prof.Sample))
	summary.Stats["totalLocations"] = int64(len(prof.Location))
	summary.Stats["totalFunctions"] = int64(len(prof.Function))

	// Extract top samples by memory usage
	if spaceIndex >= 0 && len(prof.Sample) > 0 {
		// Create a copy of samples to sort
		samples := make([]*profile.Sample, len(prof.Sample))
		copy(samples, prof.Sample)

		// Sort samples by memory usage (descending)
		sort.Slice(samples, func(i, j int) bool {
			if spaceIndex < len(samples[i].Value) && spaceIndex < len(samples[j].Value) {
				return samples[i].Value[spaceIndex] > samples[j].Value[spaceIndex]
			}
			return false
		})

		// Get top samples (limit to 10 or fewer)
		sampleLimit := 10
		if len(samples) < sampleLimit {
			sampleLimit = len(samples)
		}

		for i := 0; i < sampleLimit; i++ {
			sample := samples[i]
			if len(sample.Value) <= spaceIndex {
				continue
			}

			// Extract top functions from the stack
			var topFuncs []string
			funcLimit := 3 // Top 3 functions in the stack
			for j := 0; j < len(sample.Location) && j < funcLimit; j++ {
				loc := sample.Location[j]
				if len(loc.Line) > 0 && loc.Line[0].Function != nil {
					topFuncs = append(topFuncs, loc.Line[0].Function.Name)
				}
			}

			// Create condensed sample
			condensed := CondensedSample{
				Value:        sample.Value,
				TopFunctions: topFuncs,
			}

			// Add labels if present
			if sample.Label != nil || sample.NumLabel != nil {
				labelInfo := LabelInfo{
					Labels:   make(map[string]string),
					NumLabel: make(map[string][]int64),
				}

				for k, v := range sample.Label {
					if len(v) > 0 {
						labelInfo.Labels[k] = v[0]
					}
				}

				for k, v := range sample.NumLabel {
					labelInfo.NumLabel[k] = v
				}

				condensed.Labels = labelInfo
			}

			summary.TopSamples = append(summary.TopSamples, condensed)

			// Add to memory usage summary if we have both space and objects indices
			if spaceIndex >= 0 && objectsIndex >= 0 &&
				spaceIndex < len(sample.Value) && objectsIndex < len(sample.Value) {
				topFunc := ""
				if len(topFuncs) > 0 {
					topFunc = topFuncs[0]
				}

				summary.MemoryUsage = append(summary.MemoryUsage, MemoryUsageSummary{
					Size:  sample.Value[spaceIndex],
					Count: sample.Value[objectsIndex],
					Stack: topFunc,
				})
			}
		}
	}

	return summary, nil
}

// Find the index for memory space (either inuse_space or alloc_space)
func findMemorySpaceIndex(prof *profile.Profile) int {
	for i, sampleType := range prof.SampleType {
		if sampleType.Type == "inuse_space" || sampleType.Type == "alloc_space" {
			return i
		}
	}
	return -1
}

// Find the index for memory objects (either inuse_objects or alloc_objects)
func findMemoryObjectsIndex(prof *profile.Profile) int {
	for i, sampleType := range prof.SampleType {
		if sampleType.Type == "inuse_objects" || sampleType.Type == "alloc_objects" {
			return i
		}
	}
	return -1
}
