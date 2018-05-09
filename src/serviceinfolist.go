package gopifinder

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ServiceInfoList holds a list of Services
type ServiceInfoList struct {
	Services []ServiceInfo `json:"services"`
}

// ReadFromRequest reads the request body and deserializes it into the entity values
func (s *ServiceInfoList) ReadFromRequest(r *http.Request) error {
	if r.ContentLength != 0 {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.New("Cannot read request body. " + err.Error())
		}
		if b != nil && len(b) != 0 {
			if err := json.Unmarshal(b, &s); err != nil {
				return errors.New("Error deserializing ServiceInfoList. " + err.Error())
			}
		}
	}
	return nil
}

// WriteToResponse serializes the entity and writes it to the http response
func (s *ServiceInfoList) WriteToResponse(w http.ResponseWriter) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (s *ServiceInfoList) Serialize() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (s *ServiceInfoList) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &s)
	if err != nil {
		return err
	}
	return nil
}
