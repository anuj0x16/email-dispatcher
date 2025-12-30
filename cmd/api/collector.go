package main

import (
	"net/http"

	"github.com/anuj0x16/email-dispatcher/internal/jobs"
	"github.com/anuj0x16/email-dispatcher/internal/validator"
	"github.com/google/uuid"
)

func (app *application) emailCollector(w http.ResponseWriter, r *http.Request) {
	var input struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.To != "", "to", "must be provided")
	v.Check(validator.Matches(input.To, validator.EmailRX), "to", "must be a valid email address")
	v.Check(input.Subject != "", "subject", "must be provided")
	v.Check(len(input.Subject) <= 255, "subject", "must not be more than 255 bytes long")
	v.Check(input.Body != "", "body", "must be provided")
	v.Check(len(input.Body) <= 1_048_576, "body", "must not be more than 1,048,576 bytes long")

	if !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	job := jobs.EmailJob{
		ID:          uuid.New(),
		To:          input.To,
		Subject:     input.Subject,
		Body:        input.Body,
		Attempts:    0,
		MaxAttempts: 3,
	}

	app.dispatcher.JobQueue <- job

	data := envelope{
		"job_id": job.ID.String(),
		"status": "queued",
	}

	if err := app.writeJSON(w, http.StatusAccepted, data); err != nil {
		app.serverError(w, r, err)
	}
}
