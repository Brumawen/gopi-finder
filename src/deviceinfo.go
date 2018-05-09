package gopifinder

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
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

// AsJSON converts the current struct information to a JSON formatted string
func (d *DeviceInfo) AsJSON() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DeviceInfoFromJSON generates a DeviceInfo struct from the JSON formatted string
func DeviceInfoFromJSON(b []byte) (DeviceInfo, error) {
	var d DeviceInfo
	err := json.Unmarshal(b, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}
