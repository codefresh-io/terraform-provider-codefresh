package cfclient

import (
	"net/url"
	"strings"
)

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

func uriEncode(path string) string {
	replacer := strings.NewReplacer("+", "%20", "%2A", "*") // match Javascript's encodeURIComponent()
	return replacer.Replace(url.QueryEscape(path))
}

func UriEncodeEvent(event string) string {
	// The following is odd, but it's intentional. The event is URI encoded twice because
	// the Codefresh API expects it to be encoded twice.
	return uriEncode(uriEncode(event))
}
