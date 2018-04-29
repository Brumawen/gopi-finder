package gopifinder

// DeviceInfo holds the information about a device.
// This information details the services provided by the device.
type DeviceInfo struct {
	MachineID string
	IPAddress string
	HostName  string
	Services  []ServiceInfo
}

// GetService returns if the device provides the specified service and, if so,
// also returns the Service information
func (d *DeviceInfo) GetService(service string) (bool, *ServiceInfo) {
	for _, s := range d.Services {
		if s.Service == service {
			return true, &s
		}
	}
	return false, nil
}
