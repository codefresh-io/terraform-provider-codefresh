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
