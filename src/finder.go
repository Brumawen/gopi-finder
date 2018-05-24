package gopifinder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kardianos/service"
)

// Finder will search for and hold a list of devices available on the local network.
type Finder struct {
	PortNo         int            // Port number to attempt to connect to
	Devices        []DeviceInfo   // List of discovered devices
	VerboseLogging bool           // Switch on verbose logging
	Timeout        int            // The timeout in seconds to wait for a response from the LAN IP probe
	LastSearch     time.Time      // The date and time the last search was made
	ForceSearch    bool           // Indicates if a search must occur
	IsServer       bool           // Indicates this instance is a finder server
	MyInfo         *DeviceInfo    // The machine's device information
	Logger         service.Logger // The logger
}

// FindDevices searches the local LANs for devices.
// This will initiate a LAN wide search for each local IP address associated with
// the current device.
func (f *Finder) FindDevices() ([]DeviceInfo, error) {
	// Clear array
	f.Devices = []DeviceInfo{}
	if f.PortNo <= 0 {
		f.PortNo = 20502
	}
	if f.Timeout <= 0 {
		f.Timeout = 2
	}
	f.ForceSearch = false

	f.logDebug("Starting search...")

	ipLst, err := GetLocalIPAddresses()
	if err != nil {
		return nil, errors.New("Error getting Local IP Addresses. " + err.Error())
	}

	c := make(chan DeviceInfo)

	timeout := time.After(time.Duration(f.Timeout) * time.Second)

	// Start the goroutines looking for device on the networks
	count := 0
	for _, ip := range ipLst {
		f.logDebug("FindDevices: Searching LAN for IP Address", ip)
		scanList, err := GetPotentialAddresses(ip)
		if err != nil {
			return nil, errors.New("Error getting potential IP scan list. " + err.Error())
		}
		for _, scanIP := range scanList {
			count = count + 1
			myIP := scanIP
			go func() { c <- f.checkIfOnline(myIP) }()
		}
	}

	// Now listen for the results
	for i := 0; i < count; i++ {
		select {
		case result := <-c:
			f.AddDevice(result)
		case <-timeout:
			break
		}
	}

	f.logDebug("Completed search.")

	f.LastSearch = time.Now()
	return f.Devices, nil
}

// GetMyInfo returns the latest device information for the current device.
func (f *Finder) GetMyInfo() (DeviceInfo, bool, error) {
	if f.MyInfo == nil || len(f.MyInfo.IPAddress) == 0 || time.Since(f.MyInfo.Created).Minutes() > 5 {
		info, err := NewDeviceInfo()
		f.MyInfo = &info
		return info, true, err
	}
	return *f.MyInfo, false, nil
}

// RegisterServices registers the list of services with the
// registered devices on the network.
func (f *Finder) RegisterServices(sl []ServiceInfo) error {
	// First contact a device to get the list of devices
	devList, err := f.getCurrentDeviceList()
	if err != nil {
		return err
	}

	for _, i := range devList {
		d := i
		for n := 0; n < len(i.IPAddress); n++ {
			ln := n
			go f.registerServices(d, ln, sl)
		}
	}

	return nil
}

// SearchForDevices will search the registered devices for services that match the
// list of service names specified.
func (f *Finder) SearchForDevices() ([]DeviceInfo, error) {
	// First contact a device to get the list of devices
	f.ForceSearch = true
	devList, err := f.getCurrentDeviceList()
	if err != nil {
		return devList, err
	}
	return devList, nil
}

// SearchForServices will search the registered devices for services that match the
// list of service names specified.
func (f *Finder) SearchForServices() ([]ServiceInfo, error) {
	// First contact a device to get the list of devices
	srvList := []ServiceInfo{}
	devList, err := f.getCurrentDeviceList()
	if err != nil {
		return srvList, err
	}

	c := make(chan []ServiceInfo)
	timeout := time.After(time.Duration(f.Timeout) * time.Second)

	for _, i := range devList {
		d := i
		for n := 0; n < len(i.IPAddress); n++ {
			ln := n
			go func() { c <- f.scanForServices(d, ln) }()
		}
	}

	// Now listen for the results
	for i := 0; i < len(devList); i++ {
		select {
		case result := <-c:
			for _, r := range result {
				srvList = append(srvList, r)
			}
		case <-timeout:
			break
		}
	}

	return srvList, nil
}

func (f *Finder) getURL(ip string, method string) string {
	return fmt.Sprintf("http://%s:%d%s", ip, f.PortNo, method)
}

func (f *Finder) getCurrentDeviceList() ([]DeviceInfo, error) {
	f.logDebug("Getting current device list.")
	if len(f.Devices) == 0 {
		f.logDebug("Local list is empty.  Searching for devices.")
		return f.FindDevices()
	}
	// Check to see if we need to do a full search
	if f.ForceSearch {
		f.logDebug("Force search is set.  Searching for devices.")
		// Send a message to each of our current devices and
		// accept the device list from the first response back
		c := make(chan []DeviceInfo)
		timeout := time.After(time.Duration(f.Timeout) * time.Second)
		for _, i := range f.Devices {
			d := i
			for n := 0; n < len(i.IPAddress); n++ {
				ln := n
				go func() { c <- f.scanForDevices(d, ln) }()
			}
		}
		// Listen for the first response
		select {
		case result := <-c:
			f.Devices = []DeviceInfo{}
			for _, r := range result {
				f.Devices = append(f.Devices, r)
			}
		case <-timeout:
			f.logDebug("Search timed out.")
			break
		}
	}
	f.ForceSearch = false
	return f.Devices, nil
}

// AddDevice adds the specified device to the devices list
func (f *Finder) AddDevice(d DeviceInfo) {
	isNew := true
	if d.MachineID == "" {
		isNew = false
	} else {
		for _, i := range f.Devices {
			if i.MachineID == d.MachineID {
				isNew = false
				break
			}
		}
	}
	if isNew {
		f.Devices = append(f.Devices, d)
	}
}

func (f *Finder) checkIfOnline(ip string) DeviceInfo {
	d := DeviceInfo{}

	// Try to call the online web service of the device
	timeout := time.Duration(time.Duration(f.Timeout) * time.Second)
	client := http.Client{Timeout: timeout}
	if f.IsServer {
		// Send the current server's DeviceInfo in the call as well
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(f.MyInfo)
		if response, err := client.Post(f.getURL(ip, "/online"), "application/json;charset=utf-8", b); err == nil {
			if response.ContentLength != 0 {
				if err := d.ReadFrom(response.Body); err != nil {
					f.logError("Error reading Online Response from", ip, err.Error())
				}
			}
		}
	} else {
		if response, err := client.Get(f.getURL(ip, "/online")); err == nil {
			if response.ContentLength != 0 {
				if err := d.ReadFrom(response.Body); err != nil {
					f.logError("Error reading Online Response from", ip, err.Error())
				}
			}
		}
	}
	return d
}

func (f *Finder) scanForServices(d DeviceInfo, ipNo int) []ServiceInfo {
	client := http.Client{}
	if response, err := client.Get(d.GetURL(ipNo, "/service/get")); err != nil {
		time.Sleep(time.Duration(f.Timeout+1) * time.Second)
	} else {
		if response.ContentLength != 0 {
			siList := ServiceInfoList{}
			if err := siList.ReadFrom(response.Body); err != nil {
				f.logError("Error reading Service List response from", d.HostName, err.Error())
			} else {
				return siList.Services
			}
		}
	}
	return []ServiceInfo{}
}

func (f *Finder) scanForDevices(d DeviceInfo, ipNo int) []DeviceInfo {
	client := http.Client{}
	if response, err := client.Get(d.GetURL(ipNo, "/device/get")); err != nil {
		time.Sleep(time.Duration(f.Timeout+1) * time.Second)
	} else {
		if response.ContentLength != 0 {
			diList := DeviceInfoList{}
			if err := diList.ReadFrom(response.Body); err != nil {
				f.logError("Error reading Device List response from", d.HostName, err.Error())
			} else {
				return diList.Devices
			}
		}
	}
	return []DeviceInfo{}
}

func (f *Finder) registerServices(d DeviceInfo, ipNo int, sl []ServiceInfo) error {
	// Create a ServiceInfoList object that will be used to hold the ServiceInfo slice
	siList := ServiceInfoList{Services: sl}
	// Post the list to the device
	client := http.Client{}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(siList)
	if _, err := client.Post(d.GetURL(ipNo, "/service/add"), "application/json;charset=utf-8", b); err != nil {
		return err
	}
	return nil
}

func (f *Finder) logDebug(v ...interface{}) {
	if f.VerboseLogging {
		a := fmt.Sprint(v)
		if f.Logger == nil {
			log.Println(a[1 : len(a)-1])
		} else {
			f.Logger.Info("Finder: ", a[1:len(a)-1])
		}
	}
}

func (f *Finder) logError(v ...interface{}) {
	a := fmt.Sprint(v)
	if f.Logger == nil {
		log.Println(a[1 : len(a)-1])
	} else {
		f.Logger.Error("Finder: ", a[1:len(a)-1])
	}
}
