package util

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// this is for sourcing an env file

func LoadEnv(path string) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		l := strings.TrimSpace(s.Text())

		// continue when we encounter a # or "" or a newline
		if strings.HasPrefix(l, "#") || l == "" {
			continue
		}

		// parse the key/value and put it in the environment
		kv := strings.Split(s.Text(), "=")
		if len(kv) < 2 {
			log.Fatalf("invalid line for: %s", s.Text())
		}
		os.Setenv(kv[0], kv[1])
	}

	if err := s.Err(); err != nil {
		log.Fatalf("error reading the file: %v", err)
	}
}
