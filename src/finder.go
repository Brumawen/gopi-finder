package gopifinder

import gopitools "github.com/brumawen/gopi-tools/src"

// Finder will search for and hold a list of devices on the local network
// and the services that each device provides.
type Finder struct {
	Devices []DeviceInfo
}

func (f *Finder) FindDevices(includeMe bool) error {
	if ipLst, err := gopitools.GetLocalIPAddresses(); err != nil {
		return err
	} else {
		for _, ip := range ipLst {

		}
	}
}

// GetFirstService returns if a device in the list provides the specified service
// and, if so, also returns the Service information.
func (f *Finder) GetFirstService(service string) (bool, *ServiceInfo) {
	if f.Devices != nil {
		for _, i := range f.Devices {
			if b, s := i.GetService(service); b {
				return true, s
			}
		}
	}
	return false, nil
}
