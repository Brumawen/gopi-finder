package main

import "github.com/gorilla/mux"
import "net/http"

type ServiceController struct {
	Srv Service
}

func (c *ServiceController) AddServiceController(router *mux.Router) {
	router.Methods("POST").Path("/service/add").Name("AddService").
		Handler(Logger(http.HandlerFunc(c.handleAddService)))
	router.Methods("DELETE").Path("/service/remove/{id}/{name}").Name("RemoveService").
		Handler(Logger(http.HandlerFunc(c.handleRemoveService)))
	router.Methods("DELETE").Path("/service/remove/{id}").Name("RemoveAll").
		Handler(Logger(http.HandlerFunc(c.handleRemoveAll)))
	router.Methods("GET").Path("/service/get").Name("GetLocalServices").
		Handler(Logger(http.HandlerFunc(c.handleGetLocal)))
	router.Methods("GET").Path("/service/search/{name}").Name("Search").
		Handler(Logger(http.HandlerFunc(c.handleGetLocal)))

}

func (c *ServiceController) handleAddService(w http.ResponseWriter, r *http.Request) {

}

func (c *ServiceController) handleRemoveService(w http.ResponseWriter, r *http.Request) {

}

func (c *ServiceController) handleRemoveAll(w http.ResponseWriter, r *http.Request) {

}

func (c *ServiceController) handleGetLocal(w http.ResponseWriter, r *http.Request) {

}

func (c *ServiceController) handleSearch(w http.ResponseWriter, r *http.Request) {

}
