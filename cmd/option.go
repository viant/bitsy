package cmd

type Options struct {
	RuleURL        string `short:"r" long:"rule" description:"rule URL"`
	Validate       bool   `short:"V" long:"validate" description:"run validation"`
	BatchField     string `short:"b" long:"batch" description:"run validation"`
	SequenceField  string `short:"q" long:"seq" description:"sequence field"`
	SourceURL      string `short:"s" long:"sURL" description:"source URL"`
	DestinationURL string `short:"d" long:"dURL" description:"destination URL"`
	TimeField      string `short:"t" long:"tFileld" description:"time field"`
}
