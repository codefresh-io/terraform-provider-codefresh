package validation

import (
	"fmt"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/helper/validation/validationopts"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robfig/cron"
)

func CronExpression(opts ...validationopts.OptionSetter) schema.SchemaValidateDiagFunc {
	// Cron expression requirements: 6 fields, with ability to use descriptors (e.g. @yearly)
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	options := validationopts.NewValidationOptions().
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
