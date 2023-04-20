package postgres

import (
	"context"

	u "github.com/andreashanson/golang-postgres-test/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pg *pgxpool.Pool
}

func NewPostgresRepository(pg *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pg: pg}
}

func (pg *PostgresRepository) Save(ctx context.Context, name string) error {
	return nil
}

func (pg *PostgresRepository) Ping(ctx context.Context) error {
	return pg.pg.Ping(ctx)
}

func (pg *PostgresRepository) GetByID(ctx context.Context, id int) (u.User, error) {
	row := pg.pg.QueryRow(ctx, `
		SELECT name, lastname FROM users WHERE id = $1;
	`, id)

	var u u.User

	if err := row.Scan(&u.Name, &u.LastName); err != nil {
		return u, err
	}

	return u, nil
}
