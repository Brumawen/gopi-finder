package main

import (
	"io/ioutil"
	"net/http"

	"github.com/brumawen/gopi-finder/src"

	"github.com/gorilla/mux"
)

// OnlineController handles the Web Methods used to determine if a server is online.
type OnlineController struct {
	Srv Service
}

// AddOnlineController adds the routes associated with the controller to the router.
func (c *OnlineController) AddOnlineController(router *mux.Router) {
	router.Methods("POST").Path("/online").Name("OnlinePost").Handler(Logger(http.HandlerFunc(c.handleOnline)))
	router.Methods("GET").Path("/online").Name("OnlineGet").Handler(Logger(http.HandlerFunc(c.handleOnline)))
}

// handleOnline handles the /online web method call
func (c *OnlineController) handleOnline(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength != 0 {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Cannot read request body. "+err.Error(), 500)
			return
		}
		if b != nil && len(b) != 0 {
			// Get the DeviceInfo
			srcInfo, err := gopifinder.DeviceInfoFromJSON(b)
			if err != nil {
				http.Error(w, "Error deserializing DeviceInfo. "+err.Error(), 500)
				return
			}
			// Register this deviceinfo with the server
			c.Srv.RegisterDevice(srcInfo)
		}
	}

	// Get this server's deviceinfo
	myInfo, err := gopifinder.NewDeviceInfo()
	if err != nil {
		http.Error(w, "Error getting DeviceInfo. "+err.Error(), 500)
	} else {
		if output, err := myInfo.AsJSON(); err != nil {
			http.Error(w, "Error serializing DeviceInfo. "+err.Error(), 500)
		} else {
			w.Header().Set("content-type", "application/json")
			w.Write([]byte(output))
		}
	}
}
