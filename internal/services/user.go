package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"shuter-go/internal/dto"
	"shuter-go/internal/repositories"
	"time"
)

type PlayerRequest struct {
	PlayerID string
	Photos   []string
}

type PlayerService struct {
	repo         repositories.PlayerRepo
	requestChan  chan PlayerRequest
	workersCount int
}

func New(playerRepo repositories.PlayerRepo, workers int) *PlayerService {
	p := &PlayerService{
		repo:         playerRepo,
		requestChan:  make(chan PlayerRequest, 1000),
		workersCount: workers,
	}

	p.startWorkers()
	return p
}

func (p *PlayerService) startWorkers() {
	for i := 0; i < p.workersCount; i++ {
		go p.processIdentified()
	}
}

func (p *PlayerService) processIdentified() {
	for req := range p.requestChan {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		faceReq := dto.FaceIDRequest{
			PlayerID: req.PlayerID,
			Images:   req.Photos,
			Angles:   []string{"front", "left", "right", "back"}, // или динамически
		}

		faceResp, err := sendToFaceID(ctx, faceReq)
		if err != nil {
			log.Printf("Identification failed for player %s: %v", req.PlayerID, err)
			continue
		}

		credReq := dto.CredentialsRequest{
			PlayerID:   req.PlayerID,
			Images:     req.Photos,
			Angles:     []string{"front", "left", "right", "back"},
			Embeddings: faceResp.Embeddings, // если нужно
		}
		if err := p.repo.Create(ctx, credReq); err != nil {
			log.Printf("DB create failed for player %s: %v", req.PlayerID, err)
		}
	}
}

func (p *PlayerService) Create(ctx context.Context, req dto.CredentialsRequest) error {
	return p.repo.Create(ctx, req)
}

func sendToFaceID(ctx context.Context, req dto.FaceIDRequest) (*dto.FaceIDResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "http://face-id:5000/recognize", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("face-id service returned non-200 status")
	}

	var faceResp dto.FaceIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&faceResp); err != nil {
		return nil, err
	}
	return &faceResp, nil
}
