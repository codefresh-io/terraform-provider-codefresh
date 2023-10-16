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
	CronValidationOptions   *CronValidationOptions
	StringValidationOptions *StringValidationOptions
}

type ValidationOptionSetter func(*ValidationOptions)

func NewValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		severity:     diag.Error,
		summary:      "",
		detailFormat: "",
		CronValidationOptions: &CronValidationOptions{
			Parser: cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
		},
		StringValidationOptions: &StringValidationOptions{
			regexp2.RE2,
		},
	}
}

func WithSeverity(severity diag.Severity) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.SetSeverity(severity)
	}
}

func WithSummary(summary string) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.SetSummary(summary)
	}
}

func WithDetailFormat(detailFormat string) ValidationOptionSetter {
	return func(o *ValidationOptions) {
		o.SetDetailFormat(detailFormat)
	}
}

func (o *ValidationOptions) Apply(setters []ValidationOptionSetter) *ValidationOptions {
	for _, opt := range setters {
		opt(o)
	}
	return o
}

func (o *ValidationOptions) SetSeverity(severity diag.Severity) *ValidationOptions {
	o.severity = severity
	return o
}

func (o *ValidationOptions) SetSummary(summary string) *ValidationOptions {
	o.summary = summary
	return o
}

func (o *ValidationOptions) SetDetailFormat(detailFormat string) *ValidationOptions {
	o.detailFormat = detailFormat
	return o
}

func (o *ValidationOptions) GetSeverity() diag.Severity {
	return o.severity
}

func (o *ValidationOptions) GetSummary() string {
	return o.summary
}

func (o *ValidationOptions) GetDetailFormat() string {
	return o.detailFormat
}
