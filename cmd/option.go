package cmd

type Options struct {
	RuleURL        string `short:"r" long:"rule" description:"rule URL"`
	Validate       bool   `short:"V" long:"validate" description:"run validation"`
	BatchField     string `short:"b" long:"batch" description:"batch field"`
	SequenceField  string `short:"q" long:"seq" description:"sequence field"`
	SourceURL      string `short:"s" long:"sourceURL" description:"source URL"`
	DestinationURL string `short:"d" long:"destinationURL" description:"destination URL"`
	TimeField      string `short:"t" long:"timeField" description:"time field"`
	Concurrency int  `short:"c" long:"concurency" description:"processor concurrency"`
}


