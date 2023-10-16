package schemautil

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robfig/cron"
)

// CronValidationOptions contains options for validating cron expressions.
type CronValidationOptions struct {
	Parser cron.Parser
}

// SetParser sets the cron parser.
func (o *ValidationOptions) SetParser(parser cron.Parser) *ValidationOptions {
	o.CronValidationOptions.Parser = parser
	return o
}

// GetParser returns the cron parser.
func (o *ValidationOptions) GetParser() *cron.Parser {
	return &o.CronValidationOptions.Parser
}

// CronExpression returns a SchemaValidateDiagFunc that validates a cron expression.
func CronExpression(opts ...ValidationOptionSetter) schema.SchemaValidateDiagFunc {
	// Cron expression requirements: 6 fields, with ability to use descriptors (e.g. @yearly)
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	options := NewValidationOptions().
		SetSeverity(diag.Error).
		SetSummary("Invalid cron expression.").
		SetDetailFormat("The cron expression %q is invalid: %s").
		SetParser(parser).
		Apply(opts)

	return func(v interface{}, path cty.Path) (diags diag.Diagnostics) {
		expression := v.(string)

		if _, err := options.GetParser().Parse(expression); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: options.GetSeverity(),
				Summary:  options.GetSummary(),
				Detail:   fmt.Sprintf(options.GetDetailFormat(), expression, err),
			})
		}

		return diags
	}
}
