package validationopts

import "github.com/dlclark/regexp2"

type StringValidationOptions struct {
	RegexOptions regexp2.RegexOptions
}

func (o *ValidationOptions) SetRegexType(regexOptions regexp2.RegexOptions) *ValidationOptions {
	o.StringValidationOptions.RegexOptions = regexOptions
	return o
}

func (o *ValidationOptions) GetRegexType() regexp2.RegexOptions {
	return o.StringValidationOptions.RegexOptions
}
