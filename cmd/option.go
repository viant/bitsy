package cmd

type Options struct {
	RuleURL        string            `short:"r" long:"rule" description:"rule URL"`
	Validate       bool              `short:"V" long:"validate" description:"run validation"`
	BatchField     string            `short:"b" long:"batch" description:"batch field"`
	SequenceField  string            `short:"q" long:"seq" description:"sequence field"`
	SourceURL      string            `short:"s" long:"sourceURL" description:"source URL"`
	DestinationURL string            `short:"d" long:"destinationURL" description:"destination URL"`
	TimeField      string            `short:"t" long:"timeField" description:"time field"`
	IndexingFields map[string]string `short:"f" long:"fields" description:"indexing fields, e.g.: -f x:string -f y:int"`
	Concurrency    int               `short:"c" long:"concurrency" description:"processor concurrency"`
	Compress       bool              `short:"a" long:"archive" description:"gzip output"`
	ConfigURL      string            `short:"C" long:"configURL" description:"configuration url"`
}

