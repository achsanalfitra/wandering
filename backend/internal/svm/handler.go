package svm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

// handlers for SVM

type Handler struct {
	Database *sql.DB
	// other dependencies here
}

func (h *Handler) ValidateInput(w http.ResponseWriter, rq *http.Request) {
	var i []string
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(rq.Body).Decode(&i)
	if err != nil {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	// connect to database
	db := h.Database

	for _, k := range i {
		var id int64
		err := db.QueryRow(`SELECT id FROM cannonical_order WHERE real_vibe=$1`, k).Scan(&id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "cannot connect to database", http.StatusInternalServerError)
			return

		}
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "input OK"})
}
