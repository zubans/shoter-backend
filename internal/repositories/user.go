package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"shuter-go/internal/dto"
)

type DBPlayerRepo struct {
	db *sql.DB
}

type PlayerRepo interface {
	Create(ctx context.Context, req dto.CredentialsRequest) error
}

func NewDBPlayerRepo(db *sql.DB) *DBPlayerRepo {
	return &DBPlayerRepo{db: db}
}

func (r *DBPlayerRepo) Create(ctx context.Context, req dto.CredentialsRequest) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("transaction begin error: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `INSERT INTO players (player_id) VALUES ($1)`, req.PlayerID)
	for _, image := range req.Images {
		_, err = tx.ExecContext(ctx, `INSERT INTO photos (player_id, image) VALUES ($1, $2)`, req.PlayerID, image)
	}

	if err != nil {
		return fmt.Errorf("Create player error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %w", err)
	}

	return nil
}
