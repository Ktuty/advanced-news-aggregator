package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"news/internal/repository"
)

type API struct {
	db     *repository.Repository
	router *mux.Router
}

// конструктор для создания экземпляра API хендлера
func NewHandler(db *repository.Repository) *API {
	api := &API{
		db: db,
	}

	return api
}

// Создание роутера
func (api *API) InitRouts() *mux.Router {
	api.router = mux.NewRouter()
	api.router.Use(api.requestIDMiddleware)
	api.router.Use(api.logRequestMiddleware)
	api.endpoints()

	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/news", api.Posts).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/search", api.SearchPostsByTitle).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}
