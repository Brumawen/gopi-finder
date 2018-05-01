package gopifinder

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	gopitools "github.com/brumawen/gopi-tools/src"
)

// DeviceInfo holds the information about a device.
// This information details the services provided by the device.
type DeviceInfo struct {
	MachineID string        `json:"machineID"`
	IPAddress []string      `json:"ipAddress"`
	HostName  string        `json:"hostName"`
	Services  []ServiceInfo `json:"services"`
}

// NewDeviceInfo creates a new DeviceInfo struct and populates it with the values
// for the current device
func NewDeviceInfo() DeviceInfo {
	d := DeviceInfo{}

	//get hostname
	if out, err := exec.Command("hostname").Output(); err != nil {
		log.Println("DeviceInfo: Could not get HostName.", err)
	} else {
		d.HostName = strings.TrimSpace(string(out))
	}

	// Get the IP addresses
	if ip, err := gopitools.GetLocalIPAddresses(); err != nil {
		log.Println("DeviceInfo: Could not get IP addresses.", err)
	} else {
		d.IPAddress = ip
	}

	// Get the Machine ID
	if txt, err := gopitools.ReadAllText("/etc/machine-id"); err != nil {
		log.Println("DeviceInfo: Could not get Machine ID")
	} else {
		d.MachineID = txt
	}

	// Get the services
	d.Services = getServiceInfo()

	return d
}

func getServiceInfo() []ServiceInfo {
	fn := "serviceinfo.json"
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		log.Println("serviceinfo.json file is missing.")
		return []ServiceInfo{}
	}
	if data, err := ioutil.ReadFile(fn); err != nil {
		log.Println("DeviceInfo: Error reading serviceingo.json file.")
		return []ServiceInfo{}
	} else {
		var si ServiceInfoList
		if err := json.Unmarshal(data, &si); err != nil {
			log.Println("DeviceInfo: Error deserializing serviceinfo.json file data.")
			return []ServiceInfo{}
		} else {
			return si.Services
		}
	}
}

// GetService returns if the device provides the specified service and, if so,
// also returns the Service information
func (d *DeviceInfo) GetService(service string) (bool, *ServiceInfo) {
	for _, s := range d.Services {
		if s.ServiceName == service {
			return true, &s
		}
	}
	return false, nil
}
