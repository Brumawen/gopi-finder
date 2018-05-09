package gopifinder

// ServiceInfoList holds a list of Services
type ServiceInfoList struct {
	Services []ServiceInfo `json:"services"`
}

// ServiceInfo holds the information about a service provided by
// a device.
// The information holds the ServiceName, the Port No on the device and
// the API url stub of the service controller.
type ServiceInfo struct {
	ServiceName string `json:"serviceName"`
	MachineID   string `json:"machineID"`
	Host        string `json:"host"`
	IPAddress   string `json:"ip"`
	PortNo      int    `json:"portNo"`
	APIStub     string `json:"apiStub"`
}
