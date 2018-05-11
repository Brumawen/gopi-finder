package main

import (
	"flag"
	"fmt"

	"github.com/brumawen/gopi-finder/src"
)

func main() {
	// Subcommands
	devCmd := flag.Bool("devices", false, "The app will get a list of devices from the specified device.")
	srvCmd := flag.Bool("services", false, "The app will get a list of services from the specified device.")

	// Flag pointers
	ip := flag.String("ip", "", "IP Address of the device.")
	port := flag.Int("port", 0, "Port number of the device.")
	all := flag.Bool("a", false, "Show all device or service information.")
	verbose := flag.Bool("v", false, "Verbose logging.")
	timeout := flag.Int("t", 2, "Timeout waiting for a response from a IP probe. Defaults to 2 seconds.")

	flag.Parse()

	var d []gopifinder.DeviceInfo
	var s []gopifinder.ServiceInfo
	var err error

	f := gopifinder.Finder{VerboseLog: *verbose, Timeout: *timeout}

	if *devCmd {
		// Get devices registered with a specific device
		if i, err := gopifinder.NewDeviceInfo(); err != nil {
			print(err)
		} else {
			if *ip != "" {
				i.IPAddress = []string{*ip}
			}
			if *port > 0 {
				i.PortNo = *port
			}
			f.AddDevice(i)
		}
		d, err = f.SearchForDevices()
		if err != nil {
			print(err)
		}
	} else if *srvCmd {
		// Get services registered with a specific device
		if i, err := gopifinder.NewDeviceInfo(); err != nil {
			print(err)
		} else {
			if *ip != "" {
				i.IPAddress = []string{*ip}
			}
			if *port > 0 {
				i.PortNo = *port
			}
			f.AddDevice(i)
		}
		s, err = f.SearchForServices()
		if err != nil {
			print(err)
		}
	} else {
		// Find online Devices on the LAN
		if *verbose {
			fmt.Println("Searching for devices on the LAN...")
		}
		d, err = f.FindDevices()
		if err != nil {
			print(err)
		}
	}

	if len(d) != 0 {
		for _, i := range d {
			if *all {
				fmt.Println(i.HostName, i.MachineID, i.OS, i.IPAddress)
			} else {
				fmt.Println(i.HostName, i.IPAddress)
			}
		}
		if *verbose {
			fmt.Println("Found", len(d), "Device(s).")
		}
	}

	if len(s) != 0 {
		for _, i := range s {
			if *all {
				fmt.Println(i.MachineID, i.Host, i.IPAddress, i.PortNo, i.ServiceName, i.APIStub)
			} else {
				fmt.Println(i.MachineID, i.ServiceName, i.APIStub)
			}
		}
		if *verbose {
			fmt.Println("Found", len(s), "Service(s).")
		}
	}
}
