package gopifinder

import (
	"net/http"
)

// Entity defines a type that can be deserialized from a
// http request and can be serialized to an http response
type Entity interface {
	// ReadFromRequest reads the request body and deserializes it into the entity values
	ReadFromRequest(r *http.Request) error
	// WriteToResponse serializes the entity and writes it to the http response
	WriteToResponse(w http.ResponseWriter) error
	// Serialize serializes the entity and returns the serialized string
	Serialize() (string, error)
	// Deserialize deserializes the specified string into the entity values
	Deserialize(s string) error
}
