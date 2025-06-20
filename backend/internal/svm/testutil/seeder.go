package svm_testutil

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// this file contains utility function for testing

func SeedData(db *sql.DB) {
	// open the seed
	f, err := os.Open("C:/Users/achsanalfitra/Documents/hackathons/google-platform-hackathon/wandering/backend/internal/svm/testutil/seed/seed_v1_test.csv")
	if err != nil {
		log.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()

	// read csv
	r := csv.NewReader(f)

	// struct marker
	type cvs struct {
		O int64  `json:"vibe_order"`
		K string `json:"real_vibe"`
		V string `json:"version"`
	}

	var cva []cvs

	var h bool // heading marker
	for {
		var cvi cvs
		rc, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to parse csv line: %v", err)
		}

		if !h {
			h = true
			continue
		}

		if len(rc) != 3 {
			log.Fatalf("bad seed data")
		}

		cvi.V, cvi.K = strings.TrimSpace(rc[0]), strings.TrimSpace(rc[2])
		cvi.O, err = strconv.ParseInt(strings.TrimSpace(rc[1]), 10, 64)
		if err != nil {
			log.Fatalf("bad vibe order in seed: %v", err)
		}

		cva = append(cva, cvi)
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("failed to start transaction")
	}
	defer tx.Rollback()

	// execute
	for _, i := range cva {
		if _, err := tx.Exec(`INSERT INTO cannonical_order (real_vibe, vibe_order, version) values ($1, $2, $3)`, i.V, i.O, i.K); err != nil {
			log.Fatalf("failed to insert real vibe: %s", i.V)
		}
	}

	// finally commit
	if err := tx.Commit(); err != nil {
		log.Fatalf("error committing seed data: %v", err)
	}
}
