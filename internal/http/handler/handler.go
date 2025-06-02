package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/internal/url"
)

type HTTPErrorResponse struct {
	Error string `json:"error"`
}

type HTTPURLResponse struct {
	URL string `json:"url"`
}

type Handler struct {
	service *url.Service
	logger  *slog.Logger
}

func New(service *url.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// ShortURLHandler accepts JSON {"url":"<some_url>"} and returns {"result":"<short_url>"}
func (h *Handler) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the request is JSON
	if r.Header.Get("Content-Type") != "application/json" {
		h.WriteError(w, "Content-Type is not supported", http.StatusUnsupportedMediaType)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.WriteError(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Parse the request body
	var req struct {
		URL string `json:"url"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		h.WriteError(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if the URL is valid
	if req.URL == "" {
		h.WriteError(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Shorten the URL
	shortURL, err := h.service.ShortenURL(r.Context(), req.URL)
	if err != nil {
		h.WriteError(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	// Return the short URL
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(struct {
		Result string `json:"result"`
	}{Result: shortURL})
	if err != nil {
		h.WriteError(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// RedirectURLHandler redirects the user to the original URL
func (h *Handler) RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	originalURL, err := h.service.GetOriginalURL(r.Context(), id)
	if err != nil {
		h.WriteError(w, "Failed to get original URL", http.StatusInternalServerError)
		return
	}

	if originalURL == "" {
		h.WriteError(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		h.WriteError(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) WriteError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(HTTPErrorResponse{Error: message})
	if err != nil {
		h.logger.Error("Failed to encode response", "error", err)
	}
}
