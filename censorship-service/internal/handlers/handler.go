package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	router *mux.Router
}

func NewHandler() *API {
	return &API{router: mux.NewRouter()}
}

func (api *API) InitRouts() *mux.Router {
	api.router = mux.NewRouter()
	api.router.Use(api.requestIDMiddleware)
	api.router.Use(api.logRequestMiddleware)
	api.registerRoutes()
	return api.router
}

func (api *API) registerRoutes() {
	api.router.HandleFunc("/validate", api.ValidateComment).Methods(http.MethodPost, http.MethodOptions)
}
