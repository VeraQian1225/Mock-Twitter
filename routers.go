package main

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

//Reads from the routes slice to translate the values to httprouter.Handle
func NewRouter(routes Routes) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		var handle httprouter.Handle

		handle = route.HandlerFunc
		handle = Logger(handle)

		router.Handle(route.Method, route.Path, handle)
	}
	return router
}

// A map to store the Posts with the ID as the key acts as the storage in lieu of an actual database.
var postDB = make(map[string]*Post)

func AllRoutes() Routes {

	routes := Routes{
		Route{"PostShow", "GET", "/", PostShow},
		Route{"UserCreate", "POST", "/create", UserCreate},
		Route{"UserLogin", "POST", "/login", UserLogin},
		Route{"UserPost", "POST", "/post", UserPost},
	}
	return routes
}