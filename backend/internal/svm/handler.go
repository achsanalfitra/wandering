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

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "input OK"})
}

func (h *Handler) AddCannonicalVibe(w http.ResponseWriter, rq *http.Request) {
	type cvs struct {
		O int64  `json:"vibe_order"`
		K string `json:"real_vibe"`
		V string `json:"version"`
	}

	var cv cvs

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(rq.Body).Decode(&cv)
	if err != nil {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	// connect to db
	db := h.Database

	// check input existence
	var id int64
	err = db.QueryRow(`SELECT id FROM cannonical_order WHERE real_vibe=$1`, cv.K).Scan(&id)
	if err == nil {
		http.Error(w, "vibe already exists; use the update command", http.StatusConflict)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cannot connect to database", http.StatusInternalServerError)
		return
	}

	if _, err := db.Exec(
		`INSERT INTO cannonical_order (real_vibe, vibe_order, status, version) VALUES ($1, $2, $3, $4)`,
		cv.K, cv.O, ACTIVE, cv.V,
	); err != nil {
		http.Error(w, "failed to insert cannonical vibe pair", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "input inserted to db"})
}

func (h *Handler) DeleteCannonicalVibe(w http.ResponseWriter, rq *http.Request) {
	type ks struct {
		K string `json:"real_vibe"`
	}

	var k ks

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(rq.Body).Decode(&k)
	if err != nil {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	// connect to database
	db := h.Database

	var f bool
	err = db.QueryRow(`SELECT frozen FROM cannonical_order WHERE real_vibe=$1`, k.K).Scan(&f)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cannot connect to database", http.StatusInternalServerError)
		return
	}

	// enforce if f is true
	if f {
		http.Error(w, "cannot modify frozen vibe", http.StatusBadRequest)
		return
	}

	// delete it
	if _, err := db.Exec(`DELETE FROM cannonical_order WHERE real_vibe=$1`, k.K); err != nil {
		http.Error(w, "failed to delete cannonical vibe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "real vibe deleted"})
}
