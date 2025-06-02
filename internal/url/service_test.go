package url

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/config"
)

func TestService_ShortenURL(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T) (URLRepository, *config.Config)
		want    string
		wantErr error
	}{
		{
			name: "valid_data_should_return_short_url",
			setup: func(t *testing.T) (URLRepository, *config.Config) {
				cfg, err := config.New()
				if err != nil {
					t.Fatalf("NewConfig() error = %v", err)
				}
				ctrl := gomock.NewController(t)
				repo := NewMockURLRepository(ctrl)
				repo.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(&ShortURL{
					ID:          "1234567",
					OriginalURL: "https://www.google.com",
				}, nil)
				return repo, cfg
			},
			want: "http://localhost:8080/1234567",
		},
		{
			name: "repository_error_should_return_error",
			setup: func(t *testing.T) (URLRepository, *config.Config) {
				cfg, err := config.New()
				if err != nil {
					t.Fatalf("NewConfig() error = %v", err)
				}
				ctrl := gomock.NewController(t)
				repo := NewMockURLRepository(ctrl)
				repo.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
				return repo, cfg
			},
			wantErr: errors.New("failed to create short URL: db error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, cfg := tc.setup(t)
			service := NewService(repo, cfg)
			got, err := service.ShortenURL(context.Background(), "https://www.google.com")
			if err != nil && tc.wantErr != nil {
				if err.Error() != tc.wantErr.Error() {
					t.Errorf("ShortenURL() error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if got != tc.want {
				t.Errorf("ShortenURL() = %v, want %v", got, tc.want)
			}
		})
	}
}
