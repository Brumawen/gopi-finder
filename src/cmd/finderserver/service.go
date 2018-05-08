package main

import (
	"log"

	"github.com/brumawen/gopi-finder/src"
)

// Service defines a struct that is passed to all controllers.
type Service struct {
	Host           string
	PortNo         int
	VerboseLogging bool
	Timeout        int
	Devices        []gopifinder.DeviceInfo
	MyDevice       gopifinder.DeviceInfo
}

// RegisterDevice will add the specified DeviceInfo object to the Devices list
func (s *Service) RegisterDevice(d gopifinder.DeviceInfo) {
	if d.MachineID == s.MyDevice.MachineID {
		// This is us
		return
	}
	if s.VerboseLogging {
		log.Println("Registering device:", d.HostName, d.MachineID)
	}
	for _, i := range s.Devices {
		if i.MachineID == d.MachineID {
			// Update the Device
			i.HostName = d.HostName
			i.IPAddress = d.IPAddress
			return
		}
	}
	// Add the device
	s.Devices = append(s.Devices, d)
}
