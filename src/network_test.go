package gopifinder

import "testing"

func TestCanGetLocalIPAddresses(t *testing.T) {
	l, err := GetLocalIPAddresses()
	if err != nil {
		t.Error(err)
	}
	if len(l) == 0 {
		t.Error("No IP addresses returned.")
	}
}

func TestCanCheckIfInternetIsOnline(t *testing.T) {
	if !IsInternetOnline() {
		t.Error("Internet is not online")
	}
}
