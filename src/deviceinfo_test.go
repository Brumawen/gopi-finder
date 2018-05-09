package gopifinder

import (
	"log"
	"testing"
)

func TestCanCreateDeviceInfo(t *testing.T) {
	d, err := NewDeviceInfo()
	if err != nil {
		t.Error(err)
	}
	if d.HostName == "" {
		t.Error("HostName was not retrieved.")
	}
	if d.MachineID == "" {
		t.Error("MachineID was not retrieved.")
	}
	if len(d.IPAddress) == 0 {
		t.Error("No IP addresses were found.")
	}
	if d.OS == "" {
		t.Error("Operating System was not retrieved.")
	}
	log.Println(d)
}
