package handler

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/config"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/internal/url"
)

func TestShortURLHandler(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		contentType    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "unsupported_content_type_should_return_415",
			body:           `{"url": "https://www.google.com"}`,
			contentType:    "text/plain",
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedBody:   `{"error":"Content-Type is not supported"}` + "\n",
		},
		{
			name:           "invalid_request_should_return_400",
			body:           `bogus`,
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Failed to parse request body"}` + "\n",
		},
	}

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
			if err != nil {
				t.Fatalf("Failed to create database pool: %v", err)
			}
			defer dbpool.Close()

			repository := url.NewRepository(dbpool, false)
			service := url.NewService(repository, cfg)
			handler := New(service, logger)
			handler.ShortURLHandler(recorder, req)

			// Check the status code.
			if recorder.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, recorder.Code)
			}

			// Check the response body.
			if got := recorder.Body.String(); got != tc.expectedBody {
				t.Errorf("Expected response body %q, got %q", tc.expectedBody, got)
			}
		})
	}
}

func BenchmarkShortURLHandler_ValidRequest(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	body := `{"url": "https://www.google.com"}`
	req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	if err != nil {
		b.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	cfg, err := config.New()
	if err != nil {
		b.Fatalf("Failed to create config: %v", err)
	}
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		b.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbpool.Close()

	repository := url.NewRepository(dbpool, false)
	service := url.NewService(repository, cfg)
	handler := New(service, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recorder.Body.Reset()
		handler.ShortURLHandler(recorder, req)
	}
}
