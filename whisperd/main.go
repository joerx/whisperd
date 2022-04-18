package whisperd

import (
	"time"
)

type Shout struct {
	ID        int64     `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
