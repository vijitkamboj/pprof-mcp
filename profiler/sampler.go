package profiler

import (
	"fmt"
	"net/http"

	"github.com/google/pprof/profile"
)

type Sampler interface {
	Path() string
	Name() ProfileType
	QueryParams() map[string]string
	Summary(*profile.Profile, ProfilerDepth) (any, error)
}

func GetParsedProfile(serverAddr string, path string, queryParams map[string]string) (*profile.Profile, error) {
	// Construct the URL for the specific profile
	url := fmt.Sprintf("http://%s/%s", serverAddr, path)
	for k, v := range queryParams {
		url += fmt.Sprintf("?%s=%s", k, v)
	}

	// Fetch the profile
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s profile: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s profile, status code: %d", path, resp.StatusCode)
	}

	prof, err := profile.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s profile: %w", path, err)
	}

	return prof, nil

}
