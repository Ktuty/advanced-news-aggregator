package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"news/internal/models"
	"strconv"
)

func (api *API) Posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	requestID := r.Context().Value("X-Request-ID")
	log.Printf("Request ID: %s", requestID)

	page := getQueryParamAsInt(r, "page", 1)
	pageSize := getQueryParamAsInt(r, "pageSize", 10)
	searchStr := r.URL.Query().Get("s")

	var posts []models.Post
	var totalPages int
	var err error

	if searchStr != "" {
		posts, totalPages, err = api.db.SearchPostsByTitle(r.Context(), searchStr, page, pageSize)
	} else {
		postID := getQueryParamAsInt(r, "post_id", 0)
		posts, totalPages, err = api.db.Posts(r.Context(), page, pageSize, postID)
	}

	if err != nil {
		log.Printf("Request ID: %s, Error: %v", requestID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Posts       []models.Post `json:"posts"`
		TotalPages  int           `json:"totalPages"`
		CurrentPage int           `json:"currentPage"`
		PageSize    int           `json:"pageSize"`
	}{
		Posts:       posts,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Request ID: %s, Error: %v", requestID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HTTP метод для поиска постов по названию с постраничной навигацией
func (api *API) SearchPostsByTitle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	requestID := r.Context().Value("X-Request-ID")
	log.Printf("Request ID: %s", requestID)

	title := r.URL.Query().Get("title")
	page := getQueryParamAsInt(r, "page", 1)
	pageSize := getQueryParamAsInt(r, "pageSize", 10)

	posts, totalPages, err := api.db.SearchPostsByTitle(r.Context(), title, page, pageSize)
	if err != nil {
		log.Printf("Request ID: %s, Error: %v", requestID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Posts       []models.Post `json:"posts"`
		TotalPages  int           `json:"totalPages"`
		CurrentPage int           `json:"currentPage"`
		PageSize    int           `json:"pageSize"`
	}{
		Posts:       posts,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Request ID: %s, Error: %v", requestID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Вспомогательная функция для получения параметра запроса как целого числа с дефолтным значением
func getQueryParamAsInt(r *http.Request, paramName string, defaultValue int) int {
	paramStr := r.URL.Query().Get(paramName)
	param, err := strconv.Atoi(paramStr)
	if err != nil || param < 1 {
		return defaultValue
	}
	return param
}
