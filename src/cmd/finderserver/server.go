package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/brumawen/gopi-finder/src"
	"github.com/gorilla/mux"
)

// Server defines the Web Server.
type Server struct {
	Host           string
	PortNo         int
	VerboseLogging bool
	Timeout        int
	Devices        []gopifinder.DeviceInfo
	MyDevice       gopifinder.DeviceInfo
	Services       []gopifinder.ServiceInfo
	Router         *mux.Router
	Finder         gopifinder.Finder
}

// AddController adds the specified web service controller to the Router
func (s *Server) AddController(c Controller) {
	c.AddController(s.Router, s)
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe() error {
	if d, err := s.Finder.FindDevices(); err != nil {
		log.Print("Error finding devices.", err.Error())
	} else {
		s.Devices = d
	}
	if info, err := gopifinder.NewDeviceInfo(); err != nil {
		log.Println("Error getting Device Information.", err.Error())
	} else {
		s.MyDevice = info
		s.AddDevice(info)
	}
	return http.ListenAndServe(fmt.Sprintf("%v:%d", s.Host, s.PortNo), s.Router)
}

// AddDevice will add the specified DeviceInfo object to the Devices list
func (s *Server) AddDevice(d gopifinder.DeviceInfo) {
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
func (s *Server) RemoveDevice(id string) {
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
func (s *Server) AddService(v gopifinder.ServiceInfo) error {
	if v.MachineID == "" || v.ServiceName == "" {
		return errors.New("Missing Service ID or Name.")
	}
	for _, i := range s.Services {
		if i.MachineID == v.MachineID && i.ServiceName == v.ServiceName {
			// Update the Service
			if s.VerboseLogging {
				log.Println("Updated ServiceName", i.ServiceName, "for MachineID", i.MachineID)
			}
			i.PortNo = v.PortNo
			i.Host = v.Host
			i.IPAddress = v.IPAddress
			i.APIStub = v.APIStub
			return nil
		}
	}
	// Add the service
	if s.VerboseLogging {
		log.Println("Added ServiceName", v.ServiceName, "for MachineID", v.MachineID)
	}
	s.Services = append(s.Services, v)
	return nil
}

// RemoveService removes the service for the specified MachineID from the Services list.
func (s *Server) RemoveService(machineID string, serviceName string) {
	if machineID == "" || serviceName == "" {
		return
	}
	for n, i := range s.Services {
		if i.MachineID == machineID && i.ServiceName == serviceName {
			if s.VerboseLogging {
				log.Println("Removed ServiceName", serviceName, "for MachineID", machineID)
			}
			s.Services = append(s.Services[:n], s.Services[n+1:]...)
			return
		}
	}
}

// RemoveAllServices removes all services associated with the specified MachineID
// from the Services list
func (s *Server) RemoveAllServices(machineID string) {
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
