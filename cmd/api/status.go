package main

import "net/http"

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{"status": "available"}

	if err := app.writeJSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}
}
