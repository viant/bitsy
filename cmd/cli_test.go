package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox"
	"path"
	"testing"
)

func TestRunClient(t *testing.T) {

	parent := toolbox.CallerDirectory(3)

	var useCases = []struct {
		description string
		args        []string
		expected    int
	}{
		{
			description: " test valid yaml ",
			args: []string{
				"", "-V", "-r", path.Join(parent, "test_data/valid.yaml"),
			},
			expected: 0,
		},

		{
			description: " generate valid yaml ",
			args: []string{
				"", "-V",
				"-s", path.Join(parent, "test_data/data.json"),
				"-d", "/tmp/bitsy/$fragment/data.json",
				"-b", "batchId",
				"-q", "seq",
				"-i", "x:string,y:int",
			},
			expected: 0,
		},
		{
			description: " test invalid yaml ",
			args: []string{
				"", "-V", "-r", path.Join(parent, "test_data/invalid.yaml"),
			},
			expected: 1,
		},
		{
			description: " test im yaml ",
			args: []string{
				"", "-r", path.Join(parent, "test_data/valid.yaml"),
				"-s", path.Join(parent, "test_data/data.json"),
				"-d", "mem://localhost/tmp/bitsy/$fragment/data.json",
			},
			expected: 0,
		},
	}

	for _, useCase := range useCases {
		actual := RunClient("", useCase.args)
		assert.EqualValues(t, useCase.expected, actual, useCase.description)

	}

}
