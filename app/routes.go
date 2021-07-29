package app

import (
	"github.com/adiputra22/golangtoko/app/controllers"
	"github.com/gorilla/mux"
)

func (server *Server) initializeRoutes() {
	server.Router = mux.NewRouter()

	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
