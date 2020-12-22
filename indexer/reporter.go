package indexer

import "github.com/viant/cloudless/data/processor"

type Reporter struct {
	processor.BaseReporter
}

func NewReporter() *Reporter {
	return &Reporter{
		BaseReporter: processor.BaseReporter{
			Response: &processor.Response{Status: processor.StatusOk},
		},
	}
}
