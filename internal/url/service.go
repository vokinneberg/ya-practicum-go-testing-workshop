package url

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/config"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) ShortenURL(ctx context.Context, url string) (string, error) {
	shortURL, err := s.repo.CreateURL(ctx, &ShortURL{
		ID:          generateShortURL(7),
		OriginalURL: url,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create short URL: %w", err)
	}
	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortURL.ID), nil
}

func (s *Service) GetOriginalURL(ctx context.Context, url string) (string, error) {
	shortURL, err := s.repo.GetURL(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}
	return shortURL.OriginalURL, nil
}

// generateShortURL generates a random string alphanumeric of length 7
func generateShortURL(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
