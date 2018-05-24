package main

import (
	"fmt"
	"net/http"

	gopifinder "github.com/brumawen/gopi-finder/src"
	"github.com/gorilla/mux"
)

// ServiceController handles the Web Methods used to process service discovery.
type ServiceController struct {
	Srv *Server
}

// AddController adds the routes associated with the controller to the router.
func (c *ServiceController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("POST").Path("/service/add").Name("AddService").
		Handler(Logger(c, http.HandlerFunc(c.handleAddService)))
	router.Methods("DELETE").Path("/service/remove/{id}/{name}").Name("RemoveService").
		Handler(Logger(c, http.HandlerFunc(c.handleRemoveService)))
	router.Methods("DELETE").Path("/service/remove/{id}").Name("RemoveAll").
		Handler(Logger(c, http.HandlerFunc(c.handleRemoveAll)))
	router.Methods("GET").Path("/service/get").Name("GetLocalServices").
		Handler(Logger(c, http.HandlerFunc(c.handleGetLocal)))
	router.Methods("GET").Path("/service/search").Name("Search").
		Handler(Logger(c, http.HandlerFunc(c.handleSearch)))

}

func (c *ServiceController) handleAddService(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength != 0 {
		// Get the ServiceInfo List from the content
		l := gopifinder.ServiceInfoList{}
		if err := l.ReadFrom(r.Body); err != nil {
			http.Error(w, err.Error(), 400)
		} else {
			for _, i := range l.Services {
				// Register this service with the server
				err = c.Srv.AddService(i)
				if err != nil {
					http.Error(w, err.Error(), 400)
				}
			}
		}
	}
}

func (c *ServiceController) handleRemoveService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	name := vars["name"]
	if id == "" || name == "" {
		http.Error(w, "Invalid ID or Name", 400)
	} else {
		// Remove the service from the server list
		c.Srv.RemoveService(id, name)
	}
}

func (c *ServiceController) handleRemoveAll(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Invalid ID.", 400)
	} else {
		// Remove all services for this ID
		c.Srv.RemoveAllServices(id)
	}
}

func (c *ServiceController) handleGetLocal(w http.ResponseWriter, r *http.Request) {
	l := gopifinder.ServiceInfoList{Services: c.Srv.Services}
	if err := l.WriteTo(w); err != nil {
		http.Error(w, "Error serializing Service list. "+err.Error(), 500)
	}
}

func (c *ServiceController) handleSearch(w http.ResponseWriter, r *http.Request) {
	if s, err := c.Srv.Finder.SearchForServices(); err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		l := gopifinder.ServiceInfoList{Services: s}
		if err := l.WriteTo(w); err != nil {
			http.Error(w, "Error serializing Service list. "+err.Error(), 500)
		}
	}
}

// LogInfo is used to log information messages for this controller.
func (c *ServiceController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("ServiceController: ", a[1:len(a)-1])
}
