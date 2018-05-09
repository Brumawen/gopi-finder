package main

import "github.com/gorilla/mux"
import "net/http"

// ServiceController handles the Web Methods used to process service discovery.
type ServiceController struct {
	Srv *Server
}

// AddController adds the routes associated with the controller to the router.
func (c *ServiceController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("POST").Path("/service/add").Name("AddService").
		Handler(Logger(http.HandlerFunc(c.handleAddService)))
	router.Methods("DELETE").Path("/service/remove/{id}/{name}").Name("RemoveService").
		Handler(Logger(http.HandlerFunc(c.handleRemoveService)))
	router.Methods("DELETE").Path("/service/remove/{id}").Name("RemoveAll").
		Handler(Logger(http.HandlerFunc(c.handleRemoveAll)))
	router.Methods("GET").Path("/service/get").Name("GetLocalServices").
		Handler(Logger(http.HandlerFunc(c.handleGetLocal)))
	router.Methods("GET").Path("/service/search/{name}").Name("Search").
		Handler(Logger(http.HandlerFunc(c.handleSearch)))

}

func (c *ServiceController) handleAddService(w http.ResponseWriter, r *http.Request) {

}

func (c *ServiceController) handleRemoveService(w http.ResponseWriter, r *http.Request) {

}

func (c *ServiceController) handleRemoveAll(w http.ResponseWriter, r *http.Request) {

}

func (c *ServiceController) handleGetLocal(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Worked"))

}

func (c *ServiceController) handleSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte(vars["name"]))
}
