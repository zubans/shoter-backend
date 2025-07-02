package dto

type CredentialsRequest struct {
	Angles     []string             `json:"angles" binding:"required"`
	Images     []string             `json:"images" binding:"required"`
	PlayerID   string               `json:"playerID" binding:"required"`
	Embeddings map[string][]float64 `json:"embeddings"`
}
