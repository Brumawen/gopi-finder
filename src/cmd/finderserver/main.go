package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/brumawen/gopi-finder/src"

	"github.com/gorilla/mux"
)

func main() {
	var host = flag.String("h", "", "Host Name or IP Address.  (default All)")
	var port = flag.Int("p", 20502, "Port Number to listen on.")
	var verbose = flag.Bool("v", false, "Verbose logging.")
	var timeout = flag.Int("t", 2, "Timeout in seconds to wait for a response from a IP probe.")

	flag.Parse()

	s := Service{
		Host:           *host,
		PortNo:         *port,
		VerboseLogging: *verbose,
		Timeout:        *timeout,
	}
	if info, err := gopifinder.NewDeviceInfo(); err != nil {
		log.Println("Error getting Device Information", err.Error())
	} else {
		s.MyDevice = info
	}

	// Create the router
	router := mux.NewRouter().StrictSlash(true)
	// Add the controllers
	o := OnlineController{Srv: s}
	o.AddOnlineController(router)
	d := DeviceController{Srv: s}
	d.AddDeviceController(router)
	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%d", s.Host, s.PortNo), router))
}
