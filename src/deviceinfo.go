package gopifinder

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"

	gopitools "github.com/brumawen/gopi-tools/src"
)

// DeviceInfo holds the information about a device.
type DeviceInfo struct {
	MachineID string   `json:"machineID"`
	HostName  string   `json:"hostName"`
	IPAddress []string `json:"ipAddress"`
}

// NewDeviceInfo creates a new DeviceInfo struct and populates it with the values
// for the current device
func NewDeviceInfo() DeviceInfo {
	d := DeviceInfo{}

	// Get the Machine ID
	if txt, err := gopitools.ReadAllText("/etc/machine-id"); err != nil {
		log.Println("DeviceInfo: Could not get Machine ID")
	} else {
		d.MachineID = txt
	}

	// Get the Hostname
	if out, err := exec.Command("hostname").Output(); err != nil {
		log.Println("DeviceInfo: Could not get HostName.", err)
	} else {
		d.HostName = strings.TrimSpace(string(out))
	}

	// Get the IP addresses
	if ip, err := GetLocalIPAddresses(); err != nil {
		log.Println("DeviceInfo: Could not get IP addresses.", err)
	} else {
		d.IPAddress = ip
	}

	return d
}

func (d *DeviceInfo) AsJson() (string, error) {
	if b, err := json.Marshal(d); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func DeviceInfoFromJson(b []byte) (DeviceInfo, error) {
	var d DeviceInfo
	err := json.Unmarshal(b, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}
