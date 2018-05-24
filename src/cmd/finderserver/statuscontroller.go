package main

import (
	"fmt"
	"net/http"

	gopifinder "github.com/brumawen/gopi-finder/src"
	"github.com/gorilla/mux"
)

// StatusController handles the Web Methods used to get the status of the current device.
type StatusController struct {
	Srv *Server
}

// AddController adds the routes associated with the controller to the router.
func (c *StatusController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/status/get").Name("GetStatus").
		Handler(Logger(c, http.HandlerFunc(c.handleGetStatus)))
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

// LogInfo is used to log information messages for this controller.
func (c *StatusController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("StatusController: ", a[1:len(a)-1])
}
