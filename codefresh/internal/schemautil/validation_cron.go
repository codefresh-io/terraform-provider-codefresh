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
	parser cron.Parser
}

func (o *ValidationOptions) setCronParser(parser cron.Parser) *ValidationOptions {
	o.cronValidationOptions.parser = parser
	return o
}

// CronExpression returns a SchemaValidateDiagFunc that validates a cron expression.
func CronExpression(opts ...ValidationOptionSetter) schema.SchemaValidateDiagFunc {
	// Cron expression requirements: 5 fields, with ability to use descriptors (e.g. @yearly)
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	options := NewValidationOptions().
		setSeverity(diag.Error).
		setSummary("Invalid cron expression.").
		setDetailFormat("The cron expression %q is invalid: %s").
		setCronParser(parser).
		apply(opts)

	return func(v interface{}, path cty.Path) (diags diag.Diagnostics) {
		expression := v.(string)

		if _, err := options.cronValidationOptions.parser.Parse(expression); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: options.severity,
				Summary:  options.summary,
				Detail:   fmt.Sprintf(options.detailFormat, expression, err),
			})
		}

		return diags
	}
}
