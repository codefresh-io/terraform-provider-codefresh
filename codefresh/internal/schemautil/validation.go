package schemautil

import (
	"github.com/dlclark/regexp2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/robfig/cron"
)

type ValidationOptions struct {
	severity                diag.Severity
	summary                 string
	detailFormat            string
	cronValidationOptions   *CronValidationOptions
	stringValidationOptions *StringValidationOptions
}

type ValidationOptionSetter func(*ValidationOptions)

// NewValidationOptions returns a new ValidationOptions struct with default values.
func NewValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		severity:     diag.Error,
		summary:      "",
		detailFormat: "",
		cronValidationOptions: &CronValidationOptions{
			parser: cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
		},
		stringValidationOptions: &StringValidationOptions{
			regexp2.RE2,
		},
	}
}

// WithSeverity overrides the severity of the validation error.
func WithSeverity(severity diag.Severity) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.setSeverity(severity)
	}
}

// WithSummary overrides the summary of the validation error.
func WithSummary(summary string) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.setSummary(summary)
	}
}

// WithDetailFormat overrides the detail format string of the validation error.
//
// This string is passed to fmt.Sprintf.
// The verbs used in the format string depend on the implementation of the validation function.
func WithDetailFormat(detailFormat string) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.setDetailFormat(detailFormat)
	}
}

// WithParser overrides the cron parser used to validate cron expressions.
func WithCronParser(parser cron.Parser) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.setCronParser(parser)
	}
}

// WithRegexOptions overrides the regex options used to validate regular expressions.
func WithRegexOptions(options regexp2.RegexOptions) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.setRegexOptions(options)
	}
}

func (o *ValidationOptions) apply(setters []ValidationOptionSetter) *ValidationOptions {
	for _, opt := range setters {
		opt(o)
	}
	return o
}

func (o *ValidationOptions) setSeverity(severity diag.Severity) *ValidationOptions {
	o.severity = severity
	return o
}

func (o *ValidationOptions) setSummary(summary string) *ValidationOptions {
	o.summary = summary
	return o
}

func (o *ValidationOptions) setDetailFormat(detailFormat string) *ValidationOptions {
	o.detailFormat = detailFormat
	return o
}