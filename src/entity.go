package gopifinder

import (
	"io"
	"net/http"
)

// Entity defines a type that can be deserialized from a
// http request and can be serialized to an http response
type Entity interface {
	// ReadFrom reads the request body and deserializes it into the entity values
	ReadFrom(r io.ReadCloser) error
	// WriteTo serializes the entity and writes it to the http response
	WriteTo(w http.ResponseWriter) error
	// Serialize serializes the entity and returns the serialized string
	Serialize() (string, error)
	// Deserialize deserializes the specified string into the entity values
	Deserialize(s string) error
}
