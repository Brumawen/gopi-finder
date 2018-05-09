package gopifinder

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ServiceInfo holds the information about a service provided by
// a device.
// The information holds the ServiceName, the Port No on the device and
// the API url stub of the service controller.
type ServiceInfo struct {
	ServiceName string `json:"serviceName"`
	MachineID   string `json:"machineID"`
	Host        string `json:"host"`
	IPAddress   string `json:"ip"`
	PortNo      int    `json:"portNo"`
	APIStub     string `json:"apiStub"`
}

// ReadFromRequest reads the request body and deserializes it into the entity values
func (s *ServiceInfo) ReadFromRequest(r *http.Request) error {
	if r.ContentLength != 0 {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.New("Cannot read request body. " + err.Error())
		}
		if b != nil && len(b) != 0 {
			if err := json.Unmarshal(b, &s); err != nil {
				return errors.New("Error deserializing ServiceInfo. " + err.Error())
			}
		}
	}
	return nil
}

// WriteToResponse serializes the entity and writes it to the http response
func (s *ServiceInfo) WriteToResponse(w http.ResponseWriter) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (s *ServiceInfo) Serialize() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (s *ServiceInfo) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &s)
	if err != nil {
		return err
	}
	return nil
}
