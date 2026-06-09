package matches

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type teamInfoRepository struct {
	db *pgxpool.Pool
}

func NewTeamInfoRepository(db *pgxpool.Pool) TeamRepository {
	return &teamInfoRepository{db: db}
}

func (r *teamInfoRepository) FindByID(ctx context.Context, id string) (*TeamInfo, error) {
	query := `SELECT id, name, code, country, flag_url FROM teams WHERE id = $1`
	t := &TeamInfo{}
	err := r.db.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.Code, &t.Country, &t.FlagURL)
	if err != nil {
		return nil, fmt.Errorf("find team info by id: %w", err)
	}
	return t, nil
}
