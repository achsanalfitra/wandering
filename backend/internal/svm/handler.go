package svm

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// handlers for SVM

type Handler struct {
	Database *sql.DB
	// other dependencies here
}

func (h *Handler) ValidateInput(w http.ResponseWriter, rq *http.Request) {
	type is struct {
		I []string `json:"real_vibes"`
	}

	var i is
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(rq.Body).Decode(&i)
	if err != nil {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	// connect to database
	db := h.Database

	for _, k := range i.I {
		var id int64
		err := db.QueryRow(`SELECT id FROM cannonical_order WHERE real_vibe=$1`, k).Scan(&id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "cannot connect to database", http.StatusInternalServerError)
			return

		}

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("vibe %s not found", k), http.StatusNotFound)
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

func (h *Handler) AddBulkCannonicalVibe(w http.ResponseWriter, rq *http.Request) {
	type cvs struct {
		O int64  `json:"vibe_order"`
		K string `json:"real_vibe"`
		V string `json:"version"`
	}

	type cvbs struct {
		B []cvs `json:"bulk_vibes"`
	}

	var cvb cvbs

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(rq.Body).Decode(&cvb)
	if err != nil {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	// set max input size. Let's use 100 for now
	if len(cvb.B) > 100 {
		http.Error(w, "input exceed bulk add limit (100)", http.StatusRequestEntityTooLarge)
		return
	}

	// check own duplicate
	m := make(map[string]struct{})
	for _, i := range cvb.B {
		if _, ok := m[i.K]; ok {
			http.Error(w, fmt.Sprintf("vibe %s is dupicated", i.K), http.StatusBadRequest)
			return
		}
		m[i.K] = struct{}{}
	}

	// begin transaction
	db := h.Database
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for _, i := range cvb.B {
		var id int64
		err = db.QueryRow(`SELECT id FROM cannonical_order WHERE real_vibe=$1`, i.K).Scan(&id)
		if err == nil {
			http.Error(w, fmt.Sprintf("vibe %s exists in database", i.K), http.StatusConflict)
			return
		}

		if !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "cannot query the vibe", http.StatusInternalServerError)
			return
		}

		if _, err := tx.Exec(`INSERT INTO cannonical_order (real_vibe, vibe_order, version) VALUES ($1, $2, $3)`, i.K, i.O, i.V); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to bulk insert", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "input bulk inserted to db"})
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

func (h *Handler) ChangeStatus(w http.ResponseWriter, rq *http.Request) {
	type scs struct {
		K string `json:"real_vibe"`
		S Status `json:"status"`
	}

	var sc scs

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(rq.Body).Decode(&sc)
	if err != nil {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	if sc.S != ACTIVE && sc.S != DEPRECATED {
		http.Error(w, "unknown status", http.StatusBadRequest)
		return
	}

	// connect to db
	db := h.Database

	var s Status
	err = db.QueryRow(`SELECT status FROM cannonical_order WHERE real_vibe=$1`, sc.K).Scan(&s)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cannot connect to database", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, fmt.Sprintf("vibe %s not found", sc.K), http.StatusNotFound)
		return
	}

	if sc.S == s {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "status accepted, nothing to change"})
		return
	}

	if _, err := db.Exec(`UPDATE cannonical_order SET status=$1 WHERE real_vibe=$2`, sc.S, sc.K); err != nil {
		http.Error(w, "failed to update vibe status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("status changed for vibe %s", sc.K)})
}

func (h *Handler) Reorder(w http.ResponseWriter, rq *http.Request) {
	type ios struct {
		ID int64 `json:"id"`
		O  int64 `json:"vibe_order"`
	}

	var ioa []ios

	w.Header().Set("Content-Type", "application/json")

	// connect to db
	db := h.Database
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// get the unfrozen rows
	rows, err := tx.Query(`SELECT id, vibe_order FROM cannonical_order WHERE frozen=false`)
	if err != nil {
		http.Error(w, "failed to query unfrozen rows", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// assign the rows
	for rows.Next() {
		var ioi ios
		if err := rows.Scan(&ioi.ID, &ioi.O); err != nil {
			http.Error(w, "failed to scan row", http.StatusInternalServerError)
			return
		}
		ioa = append(ioa, ioi)
	}

	if len(ioa) == 0 {
		http.Error(w, "unfrozen rows not found", http.StatusNotFound)
		return
	}

	// get the last frozen order
	// order by vibe_order descending to get the latest data available
	var io ios
	err = tx.QueryRow(`SELECT id, vibe_order FROM cannonical_order WHERE frozen=true and vibe_order<$1 ORDER BY vibe_order DESC LIMIT 1`, ioa[0].ID).Scan(&io.ID, &io.O)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			io.O = 0
		} else {
			http.Error(w, "cannot connect to database", http.StatusInternalServerError)
			return
		}
	}

	if ioa[0].O-io.O > 1 {
		ioa[0].O -= ioa[0].O - io.O - 1
	}

	// sort the cannonical orders
	// we use bubble sort
	var s bool
	for i := range len(ioa) - 1 {
		s = false
		for j := range len(ioa) - 1 - i {
			if ioa[j].O > ioa[j+1].O {
				ioa[j], ioa[j+1] = ioa[j+1], ioa[j]
				s = true
			}
		}
		if !s {
			break
		}
	}

	// fix the gaps
	e := ioa[0].O

	for _, i := range ioa {
		if i.O != e {
			if _, err := tx.Exec(`UPDATE cannonical_order SET vibe_order=$1 WHERE id=$2`, e, i.ID); err != nil {
				http.Error(w, "unexpected error occured", http.StatusInternalServerError)
				return
			}
		}
		e++
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to reorder cannonical order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("data reordered from row %d to %d", ioa[0].ID, ioa[len(ioa)-1].ID)})
}

func (h *Handler) Freeze(w http.ResponseWriter, rq *http.Request) {
	var ia []int64

	w.Header().Set("Content-Type", "application/json")

	// begin transaction
	db := h.Database
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// get the unfrozen rows
	rows, err := tx.Query(`SELECT id FROM cannonical_order WHERE frozen=false`)
	if err != nil {
		http.Error(w, "failed to query unfrozen rows", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// assign rows
	for rows.Next() {
		var i int64
		if err := rows.Scan(&i); err != nil {
			http.Error(w, "failed to scan row", http.StatusInternalServerError)
			return
		}
		ia = append(ia, i)
	}

	if len(ia) == 0 {
		http.Error(w, "unfrozen rows not found", http.StatusNotFound)
		return
	}

	// flip the frozen rows
	for _, i := range ia {
		if _, err := tx.Exec(`UPDATE cannonical_order SET frozen=true WHERE id=ANY($1)`, i); err != nil {
			http.Error(w, "unexpected error occured", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to freeze cannonical order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("data is frozen from row %d to %d", ia[0], ia[len(ia)-1])})
}

func (h *Handler) CannonicalVibesToCSV(w http.ResponseWriter, rq *http.Request) {
	type cos struct {
		ID int64
		O  int64
		K  string
		S  Status
	}

	var coa []cos

	t := time.Now()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=cannonical_order_%s.csv", t.Format(time.RFC3339)))

	// connect to db
	db := h.Database

	rows, err := db.Query(`SELECT id, real_vibe, vibe_order, status FROM cannonical_order`)
	if err != nil {
		http.Error(w, "failed to query cannonical vibes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var coi cos
		if err := rows.Scan(&coi.ID, &coi.K, &coi.O, &coi.S); err != nil {
			http.Error(w, "failed to scan row", http.StatusInternalServerError)
			return
		}
		coa = append(coa, coi)
	}

	if len(coa) == 0 {
		http.Error(w, "no real vibes found", http.StatusNotFound)
		return
	}

	// prepare the csv
	cw := csv.NewWriter(w)

	// write header
	cw.Write([]string{"id", "real_vibe", "vibe_order", "status"})

	// write the data
	for _, coi := range coa {
		cw.Write([]string{
			strconv.FormatInt(coi.ID, 10),
			string(coi.K),
			strconv.FormatInt(coi.O, 10),
			string(coi.S),
		})
	}

	// flush the data
	cw.Flush()
	if err := cw.Error(); err != nil {
		http.Error(w, "failed to flush csv", http.StatusInternalServerError)
	}
}
