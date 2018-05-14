package gopifinder

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// ServiceInfoList holds a list of Services
type ServiceInfoList struct {
	Services []ServiceInfo `json:"services"`
}

// RegisterWith will register the Services with the specified device.
func (s *ServiceInfoList) RegisterWith(d DeviceInfo, ipNo int) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(s)
	client := http.Client{}
	_, err := client.Post(d.GetURL(ipNo, "/service/add"), "application/json;charset=utf-8", b)
	return err
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (s *ServiceInfoList) ReadFrom(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if b != nil && len(b) != 0 {
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
	}
	return nil
}

// WriteTo serializes the entity and writes it to the http response
func (s *ServiceInfoList) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
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
