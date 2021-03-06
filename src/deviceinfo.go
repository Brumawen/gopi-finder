package gopifinder

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	PortNo    int       `json:"portNo"`
	Created   time.Time `json:"created"`
}

// NewDeviceInfo creates a new DeviceInfo struct and populates it with the values
// for the current device
func NewDeviceInfo() (DeviceInfo, error) {
	d := DeviceInfo{Created: time.Now()}

	// Get the operating system
	out, err := exec.Command("uname").Output()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			d.OS = "WindowsNT"
		} else {
			return d, errors.New("Error getting device Operating System. " + err.Error())
		}
	} else {
		d.OS = strings.TrimSpace(string(out))
	}

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
func (d *DeviceInfo) CreateService(name string) ServiceInfo {
	s := ServiceInfo{
		ServiceName: name,
		MachineID:   d.MachineID,
		HostName:    d.HostName,
	}
	if len(d.IPAddress) != 0 {
		s.IPAddress = d.IPAddress[0]
	}
	return s
}

// GetURL returns the URL for the specified web method.
func (d *DeviceInfo) GetURL(idx int, method string) string {
	if d.PortNo <= 0 {
		d.PortNo = 20502
	}
	if len(d.IPAddress) < idx+1 {
		return fmt.Sprintf("http://%s:%d%s", d.HostName, d.PortNo, method)
	}
	return fmt.Sprintf("http://%s:%d%s", d.IPAddress[idx], d.PortNo, method)

}

// ReadFrom will read the request body and deserialize it into the entity values
func (d *DeviceInfo) ReadFrom(r io.ReadCloser) error {
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

// WriteTo will serialize the entity and write it to the http response
func (d *DeviceInfo) WriteTo(w http.ResponseWriter) error {
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
