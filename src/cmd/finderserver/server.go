package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brumawen/gopi-finder/src"
	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

// Server defines the Web Server.
type Server struct {
	PortNo         int                      // Port Number to listen on
	VerboseLogging bool                     // Switch on verbose logging
	Timeout        int                      // Timeout in seconds to wait for a LAN probe response
	Devices        []gopifinder.DeviceInfo  //List of registers services
	Services       []gopifinder.ServiceInfo // List of registered devices
	Finder         *gopifinder.Finder       // Finder client
	exit           chan struct{}            // Exit flag
	http           *http.Server             // HTTP server
	router         *mux.Router              // HTTP router
}

// Start is called when the service is starting
func (s *Server) Start(v service.Service) error {
	s.logInfo("Service starting")
	// Create a channel that will be used to block until the Stop signal is received
	s.exit = make(chan struct{})
	go s.run()
	return nil
}

// Stop is called when the service is stopping
func (s *Server) Stop(v service.Service) error {
	s.logInfo("Service stopping")
	// Close the channel, this will automatically release the block
	close(s.exit)
	return nil
}

// run will start up and run the service and wait for a Stop signal
func (s *Server) run() {
	if s.PortNo < 0 {
		s.PortNo = 20502
	}
	s.logInfo("Server listening on port", s.PortNo)

	// Create a router
	s.router = mux.NewRouter().StrictSlash(true)

	// Add the controllers
	s.AddController(new(OnlineController))
	s.AddController(new(DeviceController))
	s.AddController(new(ServiceController))
	s.AddController(new(StatusController))

	// Get our device info
	s.Finder = &gopifinder.Finder{
		VerboseLogging: s.VerboseLogging,
		Timeout:        s.Timeout,
		Logger:         logger,
		IsServer:       true,
	}
	if info, _, err := s.Finder.GetMyInfo(); err != nil {
		s.logError("Error getting Device Information.", err.Error())
	} else {
		info.PortNo = s.PortNo
		s.AddDevice(info)
	}

	// Tell other devices we are here
	go func() {
		s.ScanForDevices()
	}()

	// Create a HTTP server
	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.PortNo),
		Handler: s.router,
	}

	// Start the web server
	go func() {
		if err := s.http.ListenAndServe(); err != nil {
			s.logError("Error starting Web Server.", err.Error())
		}
	}()

	// Wait for an exit signal
	_ = <-s.exit

	// Shutdown
	s.http.Shutdown(nil)
}

// AddController adds the specified web service controller to the Router
func (s *Server) AddController(c Controller) {
	c.AddController(s.router, s)
}

// ScanForDevices scans the network for other devices.
func (s *Server) ScanForDevices() {
	// Get the current server device info
	s.logDebug("Scanning network for other devices.")
	isUp := false
	for !isUp {
		if info, _, err := s.Finder.GetMyInfo(); err != nil {
			s.logError("Error getting Device Information.", err.Error())
		} else {
			info.PortNo = s.PortNo
			s.AddDevice(info)

			if len(info.IPAddress) != 0 {
				if strings.HasPrefix(info.IPAddress[0], "169.254") {
					s.logInfo("Network is not DHCP capable yet.")
					time.Sleep(time.Minute)
				} else {
					// Network is up
					s.logDebug("Network is up")
					isUp = true
				}
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	}

	// Tell other devices we are here
	s.logDebug("Performing network device scan.")
	start := time.Now()
	if d, err := s.Finder.FindDevices(); err != nil {
		s.logError("Error finding devices.", err.Error())
	} else {
		for _, i := range d {
			s.AddDevice(i)
		}
	}
	s.logDebug("Network scan complete in", time.Since(start))
}

// AddDevice will add the specified DeviceInfo object to the Devices list
func (s *Server) AddDevice(d gopifinder.DeviceInfo) {
	s.logDebug("Registering device", d.HostName, d.MachineID, d.IPAddress)
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
	s.logDebug("Removing device for MachineID", id)
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
			s.logDebug("Updated ServiceName", i.ServiceName, "for MachineID", i.MachineID)
			i.PortNo = v.PortNo
			i.HostName = v.HostName
			i.IPAddress = v.IPAddress
			i.APIStub = v.APIStub
			return nil
		}
	}
	// Add the service
	s.logDebug("Added ServiceName", v.ServiceName, "for MachineID", v.MachineID)
	s.Services = append(s.Services, v)
	return nil
}

// RemoveService removes the service for the specified MachineID from the Services list.
func (s *Server) RemoveService(machineID string, serviceName string) error {
	if machineID == "" || serviceName == "" {
		return errors.New("Missing MachineID or ServiceName")
	}
	for n, i := range s.Services {
		if i.MachineID == machineID && i.ServiceName == serviceName {
			s.logDebug("Removed ServiceName", serviceName, "for MachineID", machineID)
			s.Services = append(s.Services[:n], s.Services[n+1:]...)
			return nil
		}
	}
	return nil
}

// RemoveAllServices removes all services associated with the specified MachineID
// from the Services list
func (s *Server) RemoveAllServices(machineID string) {
	if machineID == "" {
		return
	}

	s.logDebug("Removing all services for MachineID", machineID)

	n := []gopifinder.ServiceInfo{}
	for _, i := range s.Services {
		if i.MachineID != machineID {
			n = append(n, i)
		}
	}
	s.Services = n
}

func (s *Server) logDebug(v ...interface{}) {
	if s.VerboseLogging {
		a := fmt.Sprint(v)
		logger.Info("Server: ", a[1:len(a)-1])
	}
}

func (s *Server) logInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("Server: ", a[1:len(a)-1])
}

func (s *Server) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("Server: ", a[1:len(a)-1])
}
