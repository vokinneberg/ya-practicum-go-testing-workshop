package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/config"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/internal/http/handler/mocks"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/internal/url"
)

func TestShortURLHandler(t *testing.T) {
	testCases := []struct {
		name        string
		body        string
		setup       func(t *testing.T) (URLService, *config.Config)
		contentType string
		wantStatus  int
		wantBody    HTTPURLResponse
		wantError   HTTPErrorResponse
	}{
		{
			name: "valid_request_should_return_200",
			body: `{"url": "https://www.google.com"}`,
			setup: func(t *testing.T) (URLService, *config.Config) {
				cfg, err := config.New()
				if err != nil {
					t.Fatalf("Failed to create config: %v", err)
				}
				ctrl := gomock.NewController(t)
				service := mocks.NewMockURLService(ctrl)
				service.EXPECT().ShortenURL(gomock.All(), "https://www.google.com").Return("http://localhost:8080/abc123", nil)
				return service, cfg
			},
			contentType: "application/json",
			wantStatus:  http.StatusOK,
			wantBody: HTTPURLResponse{
				URL: "http://localhost:8080/abc123",
			},
		},
		{
			name: "unsupported_content_type_should_return_415",
			body: `{"url": "https://www.google.com"}`,
			setup: func(t *testing.T) (URLService, *config.Config) {
				cfg, err := config.New()
				if err != nil {
					t.Fatalf("Failed to create config: %v", err)
				}
				ctrl := gomock.NewController(t)
				service := mocks.NewMockURLService(ctrl)
				service.EXPECT().ShortenURL(gomock.Any(), "https://www.google.com").Return(gomock.Any().String(), nil).Times(0)
				return service, cfg
			},
			contentType: "text/plain",
			wantStatus:  http.StatusUnsupportedMediaType,
			wantError: HTTPErrorResponse{
				Error: "Content-Type is not supported",
			},
		},
		{
			name: "invalid_request_should_return_400",
			body: `bogus`,
			setup: func(t *testing.T) (URLService, *config.Config) {
				cfg, err := config.New()
				if err != nil {
					t.Fatalf("Failed to create config: %v", err)
				}
				ctrl := gomock.NewController(t)
				service := mocks.NewMockURLService(ctrl)
				service.EXPECT().ShortenURL(gomock.Any(), "bogus").Return(gomock.Any().String(), nil).Times(0)
				return service, cfg
			},
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
			wantError: HTTPErrorResponse{
				Error: "Failed to parse request body",
			},
		},
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

			service, _ := tc.setup(t)
			handler := New(service, logger)
			handler.ShortURLHandler(recorder, req)

			// Check the status code.
			if recorder.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d", tc.wantStatus, recorder.Code)
			}

			// Check the response body.
			var got HTTPURLResponse
			err = json.NewDecoder(recorder.Body).Decode(&got)
			if err != nil {
				t.Errorf("Failed to decode response body: %v", err)
			}

			if diff := cmp.Diff(tc.wantBody, got); diff != "" {
				t.Errorf("Response body mismatch (-want +got):\n%s", diff)
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

func FuzzShortURLHandler(f *testing.F) {
	// Add seed corpus
	f.Add(`{"url": "https://www.google.com"}`, "application/json")
	f.Add(`{"url": "https://example.com"}`, "text/plain")
	f.Add(`invalid json`, "application/json")
	f.Add(`{"url": "https://test.com"}`, "application/xml")

	cfg, err := config.New()
	if err != nil {
		f.Fatalf("Failed to create config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		f.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbpool.Close()

	repository := url.NewRepository(dbpool, false)
	service := url.NewService(repository, cfg)
	handler := New(service, logger)

	f.Fuzz(func(t *testing.T, body string, contentType string) {
		req, err := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", contentType)

		recorder := httptest.NewRecorder()
		handler.ShortURLHandler(recorder, req)

		// Basic validation - ensure we don't get a 500 error
		if recorder.Code == http.StatusInternalServerError {
			t.Errorf("Handler returned 500 for input: body=%q, contentType=%q", body, contentType)
		}

		// Ensure we get a valid response body
		if recorder.Body.Len() == 0 {
			t.Errorf("Empty response body for input: body=%q, contentType=%q", body, contentType)
		}
	})
}
