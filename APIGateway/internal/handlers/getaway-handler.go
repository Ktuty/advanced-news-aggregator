package handlers

import (
	"api-getaway/internal/models/Comment"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func (api *API) NewsShortDetailed(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	searchQuery := r.URL.Query().Get("s")
	requestID := r.Context().Value("X-Request-ID").(string)
	pageSizeQuery := r.URL.Query().Get("pageSize")

	// Перенаправление запроса на микросервис новостей
	req, err := http.NewRequest("GET", "http://localhost:8081/news?page="+strconv.Itoa(page)+"&s="+searchQuery+"&pageSize="+pageSizeQuery, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error fetching news", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (api *API) NewsFullDetailed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]
	requestID := r.Context().Value("X-Request-ID").(string)

	// Перенаправление запроса на микросервис новостей
	req, err := http.NewRequest("GET", "http://localhost:8081/news?post_id="+postID, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error fetching news", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	// Перенаправление запроса на микросервис комментариев
	commentsReq, err := http.NewRequest("GET", "http://localhost:8082/comments/"+postID, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	commentsReq.Header.Set("X-Request-ID", requestID)

	commentsResp, err := client.Do(commentsReq)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}
	defer commentsResp.Body.Close()

	commentsBody, err := io.ReadAll(commentsResp.Body)
	if err != nil {
		http.Error(w, "Error reading comments response", http.StatusInternalServerError)
		return
	}

	// Декодирование комментариев
	var comments []Comment.Comment
	if err := json.Unmarshal(commentsBody, &comments); err != nil {
		http.Error(w, "Error decoding comments", http.StatusInternalServerError)
		return
	}

	// Фильтрация комментариев через микросервис валидации
	filteredComments := filterComments(comments, requestID)

	// Объединение ответов
	response := struct {
		News     json.RawMessage   `json:"news"`
		Comments []Comment.Comment `json:"comments"`
	}{
		News:     body,
		Comments: filteredComments,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func filterComments(comments []Comment.Comment, requestID string) []Comment.Comment {
	var filteredComments []Comment.Comment
	for _, comment := range comments {
		if isValidComment(comment.Content, requestID) {
			filteredComments = append(filteredComments, comment)
		}
	}
	return filteredComments
}

func isValidComment(text, requestID string) bool {
	reqBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return false
	}

	req, err := http.NewRequest("POST", "http://localhost:8083/validate", bytes.NewReader(reqBody))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (api *API) Comment(w http.ResponseWriter, r *http.Request) {
	var cmt Comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	requestID := r.Context().Value("X-Request-ID").(string)

	// Перенаправление запроса на микросервис комментариев
	reqBody, err := json.Marshal(cmt)
	if err != nil {
		http.Error(w, "Error encoding request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8082/comments", bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error posting comment", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
