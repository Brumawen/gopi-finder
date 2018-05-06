package gopifinder

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Finder will search for and hold a list of devices on the local network
// and the services that each device provides.
type Finder struct {
	Devices []DeviceInfo

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

		f.isInitialized = true

		// Start a function that will append
		go func() {
			log.Println("Starting Checker.")
			for d := range f.deviceChan {
				if d.MachineID == "" {
					log.Println("Got close device")
					break
				}
				// Add device
				f.Devices = append(f.Devices, d)
			}
			log.Println("Ending Checker.")
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

	log.Println("FindDevices: Starting search...")
	if ipLst, err := GetLocalIPAddresses(); err != nil {
		log.Println("FindDevices: Error getting Local IP Addresses.", err)
		return nil, err
	} else {
		f.wg.Wait()
		f.wg.Add(1)

		// Clear array
		f.Devices = []DeviceInfo{}

		for _, ip := range ipLst {
			log.Println("FindDevices: Searching LAN for IP Address", ip)
			if scanList, err := GetPotentialAddresses(ip); err != nil {
				log.Println("FindDevices: Error getting potential IP scan list.", err)
			} else {
				for _, scanIp := range scanList {
					f.wg.Add(1)
					go f.pingIPAddress(scanIp)
				}
			}
		}

		// Wait for everything to complete
		f.wg.Done()
		f.wg.Wait()

		log.Println("FindDevices: Completed search.")
	}

	return f.Devices, nil
}

func (f *Finder) pingIPAddress(ip string) {
	defer f.wg.Done()

	//log.Println("Checking", ip)
	timeout := time.Duration(4 * time.Second)
	client := http.Client{Timeout: timeout}
	if response, err := client.Get("http://" + ip + ":20502/online"); err == nil {
		defer response.Body.Close()
		if contents, err := ioutil.ReadAll(response.Body); err == nil {
			log.Println("Received from", ip, "'", string(contents), "'")
			var d DeviceInfo
			log.Println(string(contents))
			if err := json.Unmarshal(contents, &d); err == nil {
				log.Println("FindDevices: Found device ", d.HostName, d.IPAddress)
				f.deviceChan <- d
			} else {
				log.Println("FindDevices: Error deserializing json string.", err)
			}
		}
	}
	//log.Println("Completed", ip)
}
