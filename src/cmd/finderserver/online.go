package main

import (
	"io/ioutil"
	"log"
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
		log.Println("Reading Body")
		if b, err := ioutil.ReadAll(r.Body); err != nil {
			log.Println("EEEEK", err)
			http.Error(w, "Cannot read request body. "+err.Error(), 500)
			return
		} else {
			log.Println("Body read OK")
			if b != nil && len(b) != 0 {
				log.Println("Getting Json")
				// Get the DeviceInfo
				if srcInfo, err := gopifinder.DeviceInfoFromJson(b); err != nil {
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
	myInfo := gopifinder.NewDeviceInfo()
	if output, err := myInfo.AsJson(); err != nil {
		http.Error(w, "Error serializing DeviceInfo. "+err.Error(), 500)
	} else {
		w.Header().Set("content-type", "application/json")
		w.Write([]byte(output))
	}
}
