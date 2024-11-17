package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type CommentRequest struct {
	Text string `json:"text"`
}

func (api *API) ValidateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if isValidComment(req.Text) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func isValidComment(text string) bool {
	forbiddenWords := []string{"qwerty", "йцукен", "ЕГЭ"}
	for _, word := range forbiddenWords {
		if strings.Contains(strings.ToLower(text), strings.ToLower(word)) {
			return false
		}
	}
	return true
}
