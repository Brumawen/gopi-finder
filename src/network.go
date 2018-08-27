package gopifinder

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

// GetClientID returns the unique client id (UUID) for this application.
func GetClientID() (string, error) {
	fn := "clientid" // File Name
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		// File does not exists, create a new uuid
		if uuid, err := uuid.NewV4(); err != nil {
			log.Println("Error creating GUID. " + err.Error())
		} else {
			uuidStr := uuid.String()
			log.Println("Created new Client ID.", uuidStr)
			err = ioutil.WriteFile(fn, []byte(uuidStr), 0666)
			if err != nil {
				return uuidStr, err
			}
			return uuidStr, nil
		}
	}
	// Read the uuid from the file
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("Failed to read the Client ID file. Attempting to recreate it.", err)
		if uuid, err := uuid.NewV4(); err != nil {
			log.Println("Error generating GUID. " + err.Error())
		} else {
			uuidStr := uuid.String()
			log.Println("Created new Client ID.", uuidStr)
			err = ioutil.WriteFile(fn, []byte(uuidStr), 0666)
			if err != nil {
				return uuidStr, err
			}
			return uuidStr, nil
		}
	}
	return string(data), nil
}

// GetLocalIPAddresses gets a list of valid IPv4 addresses for the local machine.
// These are addresses for networks that are currently up.
func GetLocalIPAddresses() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	l := []string{}
	for _, i := range ifaces {
		if i.Flags&net.FlagUp != 0 {
			adds, err := i.Addrs()
			if err != nil {
				return nil, err
			}
			for _, addr := range adds {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				// Only select valid IPv4 addresses that are not loopbacks
				if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
					l = append(l, ip.String())
				}
			}
		}
	}
	return l, nil
}

// GetPotentialAddresses gets a list of IP addresses in the same subnet as the
// specified IP address that could, potentially, host a server.
// Note that this only supports Class 3 subnets for now.
func GetPotentialAddresses(ip string) ([]string, error) {
	a := net.ParseIP(ip).To4()
	l := []string{}
	if a != nil {
		for i := 2; i < 255; i++ {
			a[3] = byte(i)
			l = append(l, a.String())
		}
	}
	return l, nil
}

// IsInternetOnline returns whether or not the machine is connected to the internet.
func IsInternetOnline() bool {
	resp, err := http.Get("http://www.msftncsi.com/ncsi.txt")
	if err != nil {
		return false
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return string(b) == "Microsoft NCSI"
}
