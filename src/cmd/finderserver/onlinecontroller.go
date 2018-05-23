package main

import (
	"net/http"

	"github.com/brumawen/gopi-finder/src"

	"github.com/gorilla/mux"
)

// OnlineController handles the Web Methods used to determine if a server is online.
type OnlineController struct {
	Srv *Server
}

// AddController adds the routes associated with the controller to the router.
func (c *OnlineController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("POST", "GET").Path("/online").Name("Online").
		Handler(Logger(http.HandlerFunc(c.handleOnline)))
}

// handleOnline handles the /online web method call
func (c *OnlineController) handleOnline(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength != 0 {
		// Get the DeviceInfo from the content
		srcInfo := gopifinder.DeviceInfo{}
		if err := srcInfo.ReadFrom(r.Body); err != nil {
			http.Error(w, err.Error(), 500)
		}
		if srcInfo.MachineID != "" {
			// Register this deviceinfo with the server
			c.Srv.AddDevice(srcInfo)
		}
	}

	// Get this server's deviceinfo
	myInfo, mustAdd, err := c.Srv.Finder.GetMyInfo()
	if err != nil {
		http.Error(w, "Error getting DeviceInfo. "+err.Error(), 500)
	} else {
		myInfo.PortNo = c.Srv.PortNo
		if mustAdd {
			c.Srv.AddDevice(myInfo)
		}
		if err := myInfo.WriteTo(w); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
