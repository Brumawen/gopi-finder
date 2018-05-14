package main

import (
	"net/http"

	gopifinder "github.com/brumawen/gopi-finder/src"
	"github.com/gorilla/mux"
)

type StatusController struct {
	Srv *Server
}

func (c *StatusController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/status/getstatus").Name("GetStatus").
		Handler(Logger(http.HandlerFunc(c.handleGetStatus)))
}

func (c *StatusController) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	s, err := gopifinder.NewDeviceStatus()
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		if err := s.WriteTo(w); err != nil {
			http.Error(w, "Error serializing Status. "+err.Error(), 500)
		}
	}
}