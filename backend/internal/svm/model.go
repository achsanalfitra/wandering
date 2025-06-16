package svm

// svm models

// cannonical_order table

type Status string

const (
	ACTIVE     Status = "active"
	DEPRECATED Status = "deprecated"
)

type CannonicalOrder struct {
	ID        int64  `json:"id"`
	VibeOrder int64  `json:"vibe_order"` // order index in the input vector
	RealVibe  string `json:"real_vibe"`  // the raw vibe name
	Version   string `json:"version"`    // model version
	Status    Status `json:"status"`     // ACTIVE or DEPRECATED
}

// known_combinations table

type KnownClassification struct {
	ID                int64  `json:"id"`
	VectorCombination []int  `json:"vector_combination"` // input vector (e.g., [0,1,1,...])
	AbstractVibe      string `json:"abstract_vibe"`      // encoded vibe like "coffee_lover"
	Version           string `json:"version"`            // model version
}
