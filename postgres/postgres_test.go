package postgres

import (
	"context"
	"testing"

	u "github.com/andreashanson/golang-postgres-test/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresRepository_Ping(t *testing.T) {
	type fields struct {
		pg *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test to ping db",
			fields: fields{
				pg: pg,
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PostgresRepository{
				pg: tt.fields.pg,
			}
			if err := pg.Ping(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepository.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepository_GetByID(t *testing.T) {
	type fields struct {
		pg *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    u.User
		wantErr bool
	}{
		{
			name: "get name that exists",
			fields: fields{
				pg: pg,
			},
			args:    args{ctx: context.Background(), id: 1},
			want:    u.User{Name: "andreas", LastName: "hansson"},
			wantErr: false,
		},
		{
			name:    "get name that does not exists",
			fields:  fields{pg: pg},
			args:    args{ctx: context.Background(), id: 2},
			want:    u.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PostgresRepository{
				pg: tt.fields.pg,
			}
			got, err := pg.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PostgresRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
