package jobs

import "github.com/google/uuid"

type EmailJob struct {
	ID          uuid.UUID
	To          string
	Subject     string
	Body        string
	Attempts    int
	MaxAttempts int
}
