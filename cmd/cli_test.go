package cmd

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox"
	"path"
	"testing"
)

func TestRunClient(t *testing.T) {

	parent := toolbox.CallerDirectory(3)

	var useCases = [] struct {
		description string
		args []string
		expected int
	}{
		{
			description:" test valid yaml ",
			args: []string {
				"" , "-V","-r" ,path.Join(parent,"test_data/valid.yaml"),
		},
		expected :0,
		},
	}


	for _, useCase := range useCases {
		actual  := RunClient("",useCase.args)
		assert.EqualValues(t,useCase.expected,actual,useCase.description)

	}


}
