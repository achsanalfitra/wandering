package svm

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateInput(t *testing.T) {
	// data
	vd := []byte(`{
		"real_vibes": ["loves_music", "loves_mountains"]
	}`)

	bd := []byte(`{
		"real_vibes": ["loves_music", "haha"]
	}`)

	// call valid
	rq := httptest.NewRequest("POST", "/", bytes.NewBuffer(vd))
	w := httptest.NewRecorder()

	HTest.ValidateInput(w, rq)
	if w.Code != http.StatusAccepted {
		t.Errorf("failed on valid data: %d", w.Code)
	}

	// call bad
	rq = httptest.NewRequest("POST", "/", bytes.NewBuffer(bd))
	w = httptest.NewRecorder()

	HTest.ValidateInput(w, rq)
	if w.Code != http.StatusNotFound {
		t.Errorf("failed on bad data: %d", w.Code)
	}
}
