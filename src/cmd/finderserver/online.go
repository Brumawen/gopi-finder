package main

import (
	"io/ioutil"
	"net/http"

	"github.com/brumawen/gopi-finder/src"

	"github.com/gorilla/mux"
)

type OnlineController struct {
	Srv Service
}

func (c *OnlineController) AddOnlineController(router *mux.Router) {
	router.Methods("POST").
		Path("/online").
		Name("OnlinePost").
		Handler(Logger(http.HandlerFunc(c.handleOnline)))
	router.Methods("GET").
		Path("/online").
		Name("OnlineGet").
		Handler(Logger(http.HandlerFunc(c.handleOnline)))
}

func (c *OnlineController) handleOnline(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength != 0 {
		if b, err := ioutil.ReadAll(r.Body); err != nil {
			http.Error(w, "Cannot read request body. "+err.Error(), 500)
			return
		} else {
			if b != nil && len(b) != 0 {
				// Get the DeviceInfo
				if srcInfo, err := gopifinder.DeviceInfoFromJSON(b); err != nil {
					http.Error(w, "Error deserializing DeviceInfo. "+err.Error(), 500)
					return
				} else {
					// Register this server
					c.Srv.RegisterServer(srcInfo)
				}
			}
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
