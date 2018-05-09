package gopifinder

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// DeviceInfo holds the information about a device.
type DeviceInfo struct {
	MachineID string    `json:"machineID"`
	HostName  string    `json:"hostName"`
	IPAddress []string  `json:"ipAddress"`
	OS        string    `json:"os"`
	Created   time.Time `json:"created"`
}

// NewDeviceInfo creates a new DeviceInfo struct and populates it with the values
// for the current device
func NewDeviceInfo() (DeviceInfo, error) {
	d := DeviceInfo{Created: time.Now()}

	// Get the operating system
	out, err := exec.Command("uname").Output()
	if err != nil {
		return d, errors.New("Error getting device Operating System. " + err.Error())
	}
	d.OS = strings.TrimSpace(string(out))

	// Get the Host Name
	out, err = exec.Command("hostname").Output()
	if err != nil {
		return d, errors.New("Error getting device HostName. " + err.Error())
	}
	d.HostName = strings.TrimSpace(string(out))

	// Get the Machine ID
	if strings.ToLower(d.OS) == "linux" {
		txt, err := ReadAllText("/etc/machine-id")
		if err != nil {
			return d, errors.New("Error getting device Machine-ID. " + err.Error())
		}
		d.MachineID = txt
	} else {
		txt, err := GetClientID()
		if err != nil {
			return d, errors.New("Error getting device Client-ID. " + err.Error())
		}
		d.MachineID = txt
	}
	if d.MachineID != "" {
		// Generate the SHA1 hash of the Machine ID
		sum := sha1.Sum([]byte(d.MachineID))
		d.MachineID = fmt.Sprintf("%x", sum)
	}

	// Get the IP addresses
	ip, err := GetLocalIPAddresses()
	if err != nil {
		return d, errors.New("Error getting device IP addresses. " + err.Error())
	}
	d.IPAddress = ip

	return d, nil
}

// CreateService creates and returns a new ServiceInfo struct object for the current device.
func (d *DeviceInfo) CreateService(serviceName string) ServiceInfo {
	return ServiceInfo{
		ServiceName: serviceName,
		MachineID:   d.MachineID,
		Host:        d.HostName,
		IPAddress:   d.IPAddress[0],
	}
}

// ReadFromRequest will read the request body and deserialize it into the entity values
func (d *DeviceInfo) ReadFromRequest(r *http.Request) error {
	if r.ContentLength != 0 {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.New("Cannot read request body. " + err.Error())
		}
		if b != nil && len(b) != 0 {
			if err := json.Unmarshal(b, &d); err != nil {
				return errors.New("Error deserializing DeviceInfo. " + err.Error())
			}
		}
	}
	return nil
}

// WriteToResponse will serialize the entity and write it to the http response
func (d *DeviceInfo) WriteToResponse(w http.ResponseWriter) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (d *DeviceInfo) Serialize() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (d *DeviceInfo) Deserialize(s string) error {
	err := json.Unmarshal([]byte(s), &d)
	if err != nil {
		return err
	}
	return nil
}
