package client

// Variable spec
type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CodefreshObject codefresh interface
type CodefreshObject interface {
	GetID() string
}

func FindInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
