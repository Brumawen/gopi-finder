package gopifinder

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// DeviceInfoList holds a list of Device Information
type DeviceInfoList struct {
	Devices []DeviceInfo `json:"devices"`
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (d *DeviceInfoList) ReadFrom(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if b != nil && len(b) != 0 {
		if err := json.Unmarshal(b, &d); err != nil {
			return err
		}
	}
	return nil
}

// WriteTo serializes the entity and writes it to the http response
func (d *DeviceInfoList) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (d *DeviceInfoList) Serialize() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (d *DeviceInfoList) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &d)
	if err != nil {
		return err
	}
	return nil
}
