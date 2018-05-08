package gopifinder

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

// DeviceInfo holds the information about a device.
type DeviceInfo struct {
	MachineID string   `json:"machineID"`
	HostName  string   `json:"hostName"`
	IPAddress []string `json:"ipAddress"`
	OS        string   `json:"os"`
}

// NewDeviceInfo creates a new DeviceInfo struct and populates it with the values
// for the current device
func NewDeviceInfo() (DeviceInfo, error) {
	d := DeviceInfo{}

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

	// Get the IP addresses
	ip, err := GetLocalIPAddresses()
	if err != nil {
		return d, errors.New("Error getting device IP addresses. " + err.Error())
	}
	d.IPAddress = ip

	return d, nil
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
