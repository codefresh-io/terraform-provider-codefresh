package validation

import (
	"fmt"
	"regexp"

	"github.com/codefresh-io/terraform-provider-codefresh/codefresh/helper/validation/validationopts"
	"github.com/dlclark/regexp2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This function has similar functionality to StringIsValidRegExp from the terraform plugin SDK.
// https://github.com/hashicorp/terraform-plugin-sdk/blob/695f0c7b92e26444786b8963e00c665f1b4ef400/helper/validation/strings.go#L225
// It has been modified to use the library https://github.com/dlclark/regexp2 instead of the standard regex golang package
// in order to support complex regular expressions including perl regex syntax.
// It has also been modified to conform to the SchemaValidateDiagFunc type instead of the deprecated SchemaValidateFunc type.
func StringIsValidRegExp(opts ...validationopts.OptionSetter) schema.SchemaValidateDiagFunc {
	options := validationopts.NewValidationOptions().
		SetSeverity(diag.Error).
		SetSummary("Invalid regular expression.").
		SetDetailFormat("%q: %s").
		Apply(opts)

	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if _, err := regexp2.Compile(value, regexp2.RE2); err != nil {
			diag := diag.Diagnostic{
				Severity: options.GetSeverity(),
				Summary:  options.GetSummary(),
				Detail:   fmt.Sprintf(options.GetDetailFormat(), p, err),
			}
			diags = append(diags, diag)
		}

		return diags
	}
}

func StringMatchesRegExp(regex string, opts ...validationopts.OptionSetter) schema.SchemaValidateDiagFunc {
	options := validationopts.NewValidationOptions().
		SetSeverity(diag.Error).
		SetSummary("Invalid value").
		SetDetailFormat("%s is invalid (must match %q)").
		Apply(opts)

	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		re := regexp.MustCompile(regex)
		if !re.MatchString(value) {
			diags = append(diags, diag.Diagnostic{
				Severity: options.GetSeverity(),
				Summary:  fmt.Sprintf(options.GetSummary(), p),
				Detail:   fmt.Sprintf(options.GetDetailFormat(), value, re.String()),
			})
		}
		return diags
	}
}
