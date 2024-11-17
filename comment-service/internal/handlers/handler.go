package handlers

import (
	"comment-service/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
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
	api.router.HandleFunc("/comments", api.AddComment).Methods(http.MethodPost)
	api.router.HandleFunc("/comments/{news_id}", api.CommentsByNewsID).Methods(http.MethodGet)
}
