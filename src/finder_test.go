package gopifinder

import "testing"
import "fmt"

func TestCanFindDevices(t *testing.T) {
	f := Finder{}
	defer f.Close()
	if d, err := f.FindDevices(); err != nil {
		t.Error(err)
	} else {
		for _, i := range d {
			fmt.Println("Found", i.HostName, i.IPAddress, i.MachineID)
		}
	}

}
