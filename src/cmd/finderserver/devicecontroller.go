package main

import (
	"net/http"

	"github.com/brumawen/gopi-finder/src"

	"github.com/gorilla/mux"
)

// DeviceController handles the Web Methods used to handle devices.
type DeviceController struct {
	Srv *Server
}

// AddController adds the routes associated with the controller to the router.
func (c *DeviceController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/device/get").Name("GetDevices").
		Handler(Logger(http.HandlerFunc(c.handleGetDevices)))
	router.Methods("DELETE").Path("/device/remove/{id}").Name("RemoveDevice").
		Handler(Logger(http.HandlerFunc(c.handleRemoveDevice)))
}

// handleGetDevices handles the /device/getdevices web method call
func (c *DeviceController) handleGetDevices(w http.ResponseWriter, r *http.Request) {
	l := gopifinder.DeviceInfoList{Devices: c.Srv.Devices}
	if err := l.WriteTo(w); err != nil {
		http.Error(w, "Error serializing Device list. "+err.Error(), 500)
	}
}

func (c *DeviceController) handleRemoveDevice(w http.ResponseWriter, r *http.Request) {

}
