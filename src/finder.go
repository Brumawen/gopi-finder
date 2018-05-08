package gopifinder

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Finder will search for and hold a list of devices on the local network
// and the services that each device provides.
type Finder struct {
	Devices    []DeviceInfo
	VerboseLog bool
	Timeout    int

	wg            sync.WaitGroup
	deviceChan    chan DeviceInfo
	isInitialized bool
}

// Init initializes the struct ready to be used.
func (f *Finder) Init() {
	if !f.isInitialized {
		// Initialize values
		f.wg = sync.WaitGroup{}
		f.deviceChan = make(chan DeviceInfo)
		f.Devices = []DeviceInfo{}

		if f.Timeout <= 0 {
			f.Timeout = 2
		}

		f.isInitialized = true

		// Start a function that will append
		go func() {
			if f.VerboseLog {
				log.Println("Starting channel hecker.")
			}
			for d := range f.deviceChan {
				if d.MachineID == "" {
					if f.VerboseLog {
						log.Println("Stopping channel checker.")
					}
					break
				}
				found := false
				for _, i := range f.Devices {
					if i.MachineID == d.MachineID {
						found = true
						break
					}
				}
				if !found {
					// Add device
					f.Devices = append(f.Devices, d)
				}
			}
			if f.VerboseLog {
				log.Println("Channel checker stopped.")
			}
		}()
	}
}

// Close cleans up the resources being held by the struct.
func (f *Finder) Close() {
	if f.isInitialized {
		f.deviceChan <- DeviceInfo{}
		f.isInitialized = false
	}
}

// FindDevices searches the local LANs for devices.
// This will initiate a LAN wide search for each local IP address associated with
// the current device.
func (f *Finder) FindDevices() ([]DeviceInfo, error) {
	if !f.isInitialized {
		f.Init()
	}

	if f.VerboseLog {
		log.Println("FindDevices: Starting search...")
	}
	ipLst, err := GetLocalIPAddresses()
	if err != nil {
		return nil, errors.New("Error getting Local IP Addresses. " + err.Error())
	}
	f.wg.Wait()
	f.wg.Add(1)

	// Clear array
	f.Devices = []DeviceInfo{}

	for _, ip := range ipLst {
		if f.VerboseLog {
			log.Println("FindDevices: Searching LAN for IP Address", ip)
		}
		scanList, err := GetPotentialAddresses(ip)
		if err != nil {
			return nil, errors.New("Error getting potential IP scan list. " + err.Error())
		}
		for _, scanIP := range scanList {
			f.wg.Add(1)
			go f.pingIPAddress(scanIP)
		}
	}

	// Wait for everything to complete
	f.wg.Done()
	f.wg.Wait()

	if f.VerboseLog {
		log.Println("FindDevices: Completed search.")
	}

	return f.Devices, nil
}

func (f *Finder) pingIPAddress(ip string) {
	defer f.wg.Done()

	if f.VerboseLog {
		log.Println("Checking", ip)
	}
	timeout := time.Duration(time.Duration(f.Timeout) * time.Second)
	client := http.Client{Timeout: timeout}
	if response, err := client.Get("http://" + ip + ":20502/online"); err == nil {
		defer response.Body.Close()
		if contents, err := ioutil.ReadAll(response.Body); err == nil {
			var d DeviceInfo
			if err := json.Unmarshal(contents, &d); err == nil {
				f.deviceChan <- d
			} else {
				log.Println("Error deserializing json string.", contents, err)
			}
		}
	}
	if f.VerboseLog {
		log.Println("Completed", ip)
	}
}
