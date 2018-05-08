package main

import (
	"flag"
	"fmt"
	"time"

	gopifinder "github.com/brumawen/gopi-finder/src"
)

func main() {
	var all = flag.Bool("a", false, "Show all device information.")
	var verbose = flag.Bool("v", false, "Verbose logging.")
	var timeout = flag.Int("t", 2, "Timeout waiting for a response from a IP probe. Defaults to 2 seconds.")
	flag.Parse()

	f := gopifinder.Finder{VerboseLog: *verbose, Timeout: *timeout}
	defer f.Close()
	start := time.Now()
	if d, err := f.FindDevices(); err != nil {
		fmt.Println(err)
	} else {
		dur := time.Since(start)
		if len(d) == 0 {
			fmt.Println("Found no devices.")
		} else {
			fmt.Println("Found", len(d), "Device(s).")
			for _, i := range d {
				if *all {
					fmt.Println(i.HostName, i.MachineID, i.OS, i.IPAddress)
				} else {
					fmt.Println(i.HostName, i.IPAddress)
				}
			}
		}
		fmt.Println("Completed in", dur.Round(time.Millisecond))
	}
}
