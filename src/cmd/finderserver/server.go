package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	Services       []gopifinder.ServiceInfo
	Router         *mux.Router
	Finder         *gopifinder.Finder
}

// AddController adds the specified web service controller to the Router
func (s *Server) AddController(c Controller) {
	c.AddController(s.Router, s)
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe() error {
	s.Finder.IsServer = true
	// Get the current server device info
	if info, _, err := s.Finder.GetMyInfo(); err != nil {
		log.Println("Error getting Device Information.", err.Error())
	} else {
		if s.Host != "" {
			info.IPAddress = []string{s.Host}
		}
		info.PortNo = s.PortNo
		s.AddDevice(info)
	}

	// Tell other devices we are here
	go func() {
		s.ScanForDevices()
	}()

	// Start the web server
	log.Println("Server listening on port", s.PortNo)
	return http.ListenAndServe(fmt.Sprintf("%v:%d", s.Host, s.PortNo), s.Router)
}

// ScanForDevices scans the network for other devices.
func (s *Server) ScanForDevices() {
	// Get the current server device info
	if s.VerboseLogging {
		log.Println("Scanning network for other devices.")
	}
	isUp := false
	for !isUp {
		if info, _, err := s.Finder.GetMyInfo(); err != nil {
			log.Println("Error getting Device Information.", err.Error())
		} else {
			if s.Host != "" {
				info.IPAddress = []string{s.Host}
			}
			info.PortNo = s.PortNo
			s.AddDevice(info)

			if len(info.IPAddress) != 0 {
				if strings.HasPrefix(info.IPAddress[0], "169.254") {
					log.Println("Network is not DHCP capable yet.")
					time.Sleep(time.Minute)
				} else {
					// Network is up
					if s.VerboseLogging {
						log.Println("Network is up")
					}
					isUp = true
				}
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}

	// Tell other devices we are here
	if s.VerboseLogging {
		log.Println("Performing scan.")
	}
	if d, err := s.Finder.FindDevices(); err != nil {
		log.Print("Error finding devices.", err.Error())
	} else {
		for _, i := range d {
			s.AddDevice(i)
		}
	}
	if s.VerboseLogging {
		log.Println("Scan complete.")
	}
}

// AddDevice will add the specified DeviceInfo object to the Devices list
func (s *Server) AddDevice(d gopifinder.DeviceInfo) {
	if s.VerboseLogging {
		log.Println("Registering device:", d.HostName, d.MachineID, d.IPAddress)
	}
	for _, i := range s.Devices {
		if i.MachineID == d.MachineID {
			// Update the Device
			i.HostName = d.HostName
			i.IPAddress = d.IPAddress
			i.Created = d.Created
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
		return errors.New("Missing Service ID or Name")
	}
	for _, i := range s.Services {
		if i.MachineID == v.MachineID && i.ServiceName == v.ServiceName {
			// Update the Service
			if s.VerboseLogging {
				log.Println("Updated ServiceName", i.ServiceName, "for MachineID", i.MachineID)
			}
			i.PortNo = v.PortNo
			i.HostName = v.HostName
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
