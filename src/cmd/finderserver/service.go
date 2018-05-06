package main

import (
	"github.com/brumawen/gopi-finder/src"
)

type Service struct {
}

func (s *Service) RegisterServer(d gopifinder.DeviceInfo) error {
	return nil
}
