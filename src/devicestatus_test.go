package gopifinder

import (
	"fmt"
	"testing"
)

func TestCanGetStatus(t *testing.T) {
	s, err := NewDeviceStatus()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}
