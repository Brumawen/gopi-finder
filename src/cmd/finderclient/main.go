package main

import (
	"flag"
	"fmt"
	"time"

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

	start := time.Now()

	f := gopifinder.Finder{
		VerboseLogging: *verbose,
		Timeout:        *timeout,
	}

	if *devCmd {
		// Get devices registered with a specific device
		if i, err := gopifinder.NewDeviceInfo(); err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}
	} else if *srvCmd {
		// Get services registered with a specific device
		if i, err := gopifinder.NewDeviceInfo(); err != nil {
			fmt.Println(err)
		} else {
			if *ip == "" {
				i.IPAddress = []string{"localhost"}
			} else {
				i.IPAddress = []string{*ip}
			}
			if *port > 0 {
				i.PortNo = *port
			}
			f.AddDevice(i)
		}
		s, err = f.SearchForServices()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// Find online Devices on the LAN
		if *verbose {
			fmt.Println("Searching for devices on the LAN...")
		}
		d, err = f.FindDevices()
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(d) != 0 {
		for _, i := range d {
			if *all {
				fmt.Printf("%s\t%s\t%s\t%s\n", i.HostName, i.MachineID, i.OS, i.IPAddress)
			} else {
				fmt.Printf("%s\t%s\n", i.HostName, i.IPAddress)
			}
		}
		if *verbose {
			fmt.Println("Found", len(d), "Device(s).")
		}
	}

	if len(s) != 0 {
		for _, i := range s {
			if *all {
				fmt.Printf("%s\t%s\t%s\t%d\t%s\t%s\n", i.HostName, i.ServiceName, i.IPAddress, i.PortNo, i.APIStub, i.MachineID)
			} else {
				fmt.Printf("%s\t%s\t%s\t%d\t%s\n", i.HostName, i.ServiceName, i.IPAddress, i.PortNo, i.APIStub)
			}
		}
		if *verbose {
			fmt.Println("Found", len(s), "Service(s).")
		}
	}

	if *verbose {
		fmt.Println("Completed in", time.Since(start).Seconds(), "sec")
	}
}
