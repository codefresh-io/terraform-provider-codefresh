package validation

const (
	// https://github.com/codefresh-io/hermes/blob/6d75b347cb8ff471ce970a766b2285788e5e19fe/pkg/backend/dev_compose_types.json#L226
	ValidCronMessageRegex string = `^[a-zA-Z0-9_+\s-#?.:]{2,128}$`
)
