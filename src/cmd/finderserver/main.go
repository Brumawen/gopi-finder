package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	s := Service{}
	o := OnlineController{Srv: s}

	router := mux.NewRouter().StrictSlash(true)
	o.AddOnlineController(router)
	log.Fatal(http.ListenAndServe(":20502", router))
}
