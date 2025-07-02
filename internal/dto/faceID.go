package dto

type FaceIDRequest struct {
	PlayerID string   `json:"playerId"`
	Images   []string `json:"images"`
	Angles   []string `json:"angles"`
}

type FaceIDResponse struct {
	PlayerID   string               `json:"playerId"`
	Embeddings map[string][]float64 `json:"embeddings"`
	MatchScore float64              `json:"matchScore,omitempty"`
}
