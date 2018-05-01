package gopifinder

import gopitools "github.com/brumawen/gopi-tools/src"
import "log"
import "sync"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "time"

// Finder will search for and hold a list of devices on the local network
// and the services that each device provides.
type Finder struct {
	Devices []DeviceInfo
}

// FindDevices searches the local LANs for devices.
// This will initiate a LAN wide search for each local IP address associated with
// the current device.
func (f *Finder) FindDevices(includeMe bool) error {
	log.Println("FindDevices: Starting search...")
	if ipLst, err := gopitools.GetLocalIPAddresses(); err != nil {
		log.Println("FindDevices: Error getting Local IP Addresses.", err)
		return err
	} else {
		myDevices := []DeviceInfo{}
		devices := make(chan DeviceInfo)
		wg := sync.WaitGroup{}
		for _, ip := range ipLst {
			log.Println("FindDevices: Searching LAN for IP Address", ip)
			if includeMe {
				// Get the values for mwa
				wg.Add(1)
				go f.pingIPAddress(ip, &wg, devices)
			}
			if scanList, err := gopitools.GetPotentialAddresses(ip); err != nil {
				log.Println("FindDevices: Error getting potential IP scan list.", err)
			} else {
				//scanList = scanList[:15]
				for _, scanIp := range scanList {
					wg.Add(1)
					go f.pingIPAddress(scanIp, &wg, devices)
				}
			}
		}

		go func() {
			for d := range devices {
				myDevices = append(myDevices, d)
			}
		}()

		// Wait for everything to complete
		wg.Wait()

		f.Devices = myDevices
		log.Println("FindDevices: Completed search.")
	}

	return nil
}

// GetFirstService returns if a device in the list provides the specified service
// and, if so, also returns the Service information.
func (f *Finder) GetFirstService(service string) (bool, *ServiceInfo) {
	if f.Devices == nil {
		err := f.FindDevices(true)
		if err != nil {
			return false, nil
		}
	}
	if f.Devices != nil {
		for _, i := range f.Devices {
			if b, s := i.GetService(service); b {
				return true, s
			}
		}
	}
	return false, nil
}

func (f *Finder) pingIPAddress(ip string, wg *sync.WaitGroup, c chan DeviceInfo) {
	defer wg.Done()

	//log.Println("Checking", ip)
	timeout := time.Duration(4 * time.Second)
	client := http.Client{Timeout: timeout}
	if response, err := client.Get("http://" + ip + ":20502/GetDeviceInfo"); err == nil {
		defer response.Body.Close()
		if contents, err := ioutil.ReadAll(response.Body); err == nil {
			var d DeviceInfo
			if err := json.Unmarshal(contents, &d); err == nil {
				log.Println("FindDevices: Found device ", d.HostName, d.IPAddress)
				c <- d
			} else {
				log.Println("FindDevices: Error deserializing json string.", err)
			}
		}
	}
	//log.Println("Completed", ip)
}
