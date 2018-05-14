package main

import (
	"flag"
	"log"

	"github.com/brumawen/gopi-finder/src"

	"github.com/gorilla/mux"
)

func main() {
	var host = flag.String("h", "", "Host Name or IP Address.  (default All)")
	var port = flag.Int("p", 20502, "Port Number to listen on.")
	var verbose = flag.Bool("v", false, "Verbose logging.")
	var timeout = flag.Int("t", 5, "Timeout in seconds to wait for a response from a IP probe.")

	flag.Parse()

	// Create a new server
	s := Server{
		Host:           *host,
		PortNo:         *port,
		VerboseLogging: *verbose,
		Timeout:        *timeout,
		Router:         mux.NewRouter().StrictSlash(true),
		Finder:         gopifinder.Finder{VerboseLog: *verbose, Timeout: *timeout},
	}

	// Add the controllers
	s.AddController(new(OnlineController))
	s.AddController(new(DeviceController))
	s.AddController(new(ServiceController))

	// Start the server
	log.Fatal(s.ListenAndServe())
}
