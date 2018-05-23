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
	router.Methods("GET").Path("/device/refresh").Name("RefreshDevices").
		Handler(Logger(http.HandlerFunc(c.handleRefreshDevices)))
}

// handleGetDevices handles the /device/getdevices web method call
func (c *DeviceController) handleGetDevices(w http.ResponseWriter, r *http.Request) {
	l := gopifinder.DeviceInfoList{Devices: c.Srv.Devices}
	if err := l.WriteTo(w); err != nil {
		http.Error(w, "Error serializing Device list. "+err.Error(), 500)
	}
}

func (c *DeviceController) handleRemoveDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid ID", 400)
	} else {
		// Remove the device from the list
		c.Srv.RemoveDevice(id)
	}
}

func (c *DeviceController) handleRefreshDevices(w http.ResponseWriter, r *http.Request) {
	go c.Srv.ScanForDevices()
	w.Write([]byte("Refresh Started."))
}
