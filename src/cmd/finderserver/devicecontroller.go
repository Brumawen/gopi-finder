package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// DeviceController handles the Web Methods used to handle devices.
type DeviceController struct {
	Srv Service
}

// AddDeviceController adds the routes associated with the controller to the router.
func (c *DeviceController) AddDeviceController(router *mux.Router) {
	router.Methods("GET").Path("/device/get").Name("GetDevices").
		Handler(Logger(http.HandlerFunc(c.handleGetDevices)))
	router.Methods("DELETE").Path("/device/remove/{id}").Name("RemoveDevice").
		Handler(Logger(http.HandlerFunc(c.handleRemoveDevice)))
}

// handleGetDevices handles the /device/getdevices web method call
func (c *DeviceController) handleGetDevices(w http.ResponseWriter, r *http.Request) {
	if output, err := json.Marshal(c.Srv.Devices); err != nil {
		http.Error(w, "Error serializing Device list. "+err.Error(), 500)
	} else {
		w.Header().Set("content-type", "application/json")
		w.Write(output)
	}
}

func (c *DeviceController) handleRemoveDevice(w http.ResponseWriter, r *http.Request) {

}
