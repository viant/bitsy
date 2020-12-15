package resource

import "testing"

func TestTracker_HasChanged(t *testing.T) {

	var useCases = []struct {
		description string
		baseURL        string
		resourceURL map[string]string
		expectedURL string
	}
