package url

import "time"

type ShortURL struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
}
