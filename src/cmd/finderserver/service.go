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
	Services       []gopifinder.ServiceInfo
}

// AddDevice will add the specified DeviceInfo object to the Devices list
func (s *Service) AddDevice(d gopifinder.DeviceInfo) {
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

// RemoveDevice removes the device with the specified ID from the Devices list.
func (s *Service) RemoveDevice(id string) {
	if id == "" {
		return
	}
	if s.VerboseLogging {
		log.Println("Removing device for MachineID", id)
	}
	for n, i := range s.Devices {
		if i.MachineID == id {
			s.Devices = append(s.Devices[:n], s.Devices[n+1:]...)
			return
		}
	}
	s.RemoveAllServices(id)
}

// AddService adds the specified ServiceInfo object to the Service list
func (s *Service) AddService(v gopifinder.ServiceInfo) {
	for _, i := range s.Services {
		if i.MachineID == v.MachineID && i.ServiceName == v.ServiceName {
			// Update the Service
			i.PortNo = v.PortNo
			i.Host = v.Host
			i.IPAddress = v.IPAddress
			i.APIStub = v.APIStub
			return
		}
		// Add the service
		s.Services = append(s.Services, v)
	}
}

// RemoveDevice removes the service for the specified MachineID from the Services list.
func (s *Service) RemoveService(machineID string, serviceName string) {
	if machineID == "" || serviceName == "" {
		return
	}
	if s.VerboseLogging {
		log.Println("Removing ServiceName", serviceName, "for MachineID", machineID)
	}
	for n, i := range s.Services {
		if i.MachineID == machineID && i.ServiceName == serviceName {
			s.Services = append(s.Services[:n], s.Services[n+1:]...)
			return
		}
	}
}

// RemoveAllServices removes all services associated with the specified MachineID
// from the Services list
func (s *Service) RemoveAllServices(machineID string) {
	if machineID == "" {
		return
	}
	if s.VerboseLogging {
		log.Println("Removing all services for MachineID", machineID)
	}
	n := []gopifinder.ServiceInfo{}
	for _, i := range s.Services {
		if i.MachineID != machineID {
			n = append(n, i)
		}
	}
	s.Services = n
}
