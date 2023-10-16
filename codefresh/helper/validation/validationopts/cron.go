package validationopts

import "github.com/robfig/cron"

type CronValidationOptions struct {
	Parser cron.Parser
}

func (o *ValidationOptions) SetParser(parser cron.Parser) *ValidationOptions {
	o.CronValidationOptions.Parser = parser
	return o
}

func (o *ValidationOptions) GetParser() *cron.Parser {
	return &o.CronValidationOptions.Parser
}
