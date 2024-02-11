package schemautil

import (
	"fmt"
	"regexp"

	"github.com/dlclark/regexp2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StringValidationOptions contains options for validating strings.
type StringValidationOptions struct {
	regexOptions regexp2.RegexOptions
}

// See: https://github.com/dlclark/regexp2/blob/03d34d8ad254ae4e2fb4f58e0723420efa1c7c07/regexp.go#L124-L142
func (o *ValidationOptions) setRegexOptions(regexOptions regexp2.RegexOptions) *ValidationOptions {
	o.stringValidationOptions.regexOptions = regexOptions
	return o
}

// StringIsValidRegExp returns a SchemaValidateDiagFunc which validates that a string is a valid regular expression.
//
// This function has similar functionality to StringIsValidRegExp from the terraform plugin SDK.
// https://github.com/hashicorp/terraform-plugin-sdk/blob/695f0c7b92e26444786b8963e00c665f1b4ef400/helper/validation/strings.go#L225
// It has been modified to use the library https://github.com/dlclark/regexp2 instead of the standard regex golang package
// in order to support complex regular expressions including perl regex syntax.
//
// It has also been modified to conform to the SchemaValidateDiagFunc type instead of the deprecated SchemaValidateFunc type.
func StringIsValidRegExp(opts ...ValidationOptionSetter) schema.SchemaValidateDiagFunc {
	options := NewValidationOptions().
		setSeverity(diag.Error).
		setSummary("Invalid regular expression.").
		setDetailFormat("%q: %s").
		apply(opts)

	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if _, err := regexp2.Compile(value, regexp2.RE2); err != nil {
			diag := diag.Diagnostic{
				Severity: options.severity,
				Summary:  options.summary,
				Detail:   fmt.Sprintf(options.detailFormat, p, err),
			}
			diags = append(diags, diag)
		}

		return diags
	}
}

// StringIsValidYaml returns a SchemaValidateDiagFunc which validates that a string is valid YAML.
func StringIsValidYaml(opts ...ValidationOptionSetter) schema.SchemaValidateDiagFunc {
	options := NewValidationOptions().
		setSeverity(diag.Error).
		setSummary("Invalid YAML").
		setDetailFormat("%s is not valid YAML: %s").
		apply(opts)

	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if _, err := NormalizeYamlString(value); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: options.severity,
				Summary:  options.summary,
				Detail:   fmt.Sprintf(options.detailFormat, p, err),
			})
		}
		return diags
	}
}

// StringMatchesRegExp returns a SchemaValidateDiagFunc which validates that a string matches a regular expression.
func StringMatchesRegExp(regex string, opts ...ValidationOptionSetter) schema.SchemaValidateDiagFunc {
	options := NewValidationOptions().
		setSeverity(diag.Error).
		setSummary("Invalid value").
		setDetailFormat("%s is invalid (must match %q)").
		setRegexOptions(regexp2.RE2).
		apply(opts)

	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		re := regexp.MustCompile(regex)
		if !re.MatchString(value) {
			diags = append(diags, diag.Diagnostic{
				Severity: options.severity,
				Summary:  options.summary,
				Detail:   fmt.Sprintf(options.detailFormat, value, re.String()),
			})
		}
		return diags
	}
}