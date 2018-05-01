package gopifinder

import "testing"

func TestCanFindDevices(t *testing.T) {
	f := Finder{}
	f.FindDevices(true)
}
