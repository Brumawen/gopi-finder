package main

import (
	"github.com/brumawen/gopi-finder/src"
)

// Service defines a struct that is passed to all controllers.
type Service struct {
	Host           string
	PortNo         int
	VerboseLogging bool
	Timeout        int
}

func (s *Service) RegisterServer(d gopifinder.DeviceInfo) error {
	return nil
}
