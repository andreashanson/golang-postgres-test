package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
)

var conn *sql.DB
var pg *pgxpool.Pool
var repo *PostgresRepository
var connectionString string

var (
	user     = "andreas"
	password = "abc"
	dbname   = "arne"
	host     = "localhost"
)

func startPostgresPool() *dockertest.Resource {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Failed creating docker pool: %v", err)
	}

	opts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=abc",
			"POSTGRES_USER=andreas",
			"POSTGRES_DB=arne",
		},
	}

	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		log.Fatalf("Failed running Postgres container: %v", err)
	}
	resource.Expire(180)
	postgresPort := resource.GetPort("5432/tcp")

	if err = pool.Retry(func() error {
		connectionString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, postgresPort, dbname)
		pg, err = pgxpool.New(context.Background(), connectionString)
		if err != nil {
			log.Fatalf("failed to open postgres: %v", err)
		}

		return pg.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Failed connecting to postgres container: %v", err)
	}

	return resource
}

func startContainer() *dockertest.Resource {
	ctx := context.Background()
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + user,
			"POSTGRES_DB=" + dbname,
		},
	}

	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	port := resource.GetPort("5432/tcp")

	if err = pool.Retry(func() error {
		dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

		pg, err = pgxpool.New(ctx, dsn)
		if err != nil {
			log.Fatalf("failed to create new pool error %s", err.Error())
		}

		return pg.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	if err := repo.Ping(ctx); err != nil {
		log.Fatalf("failed to ping pool %v", err)
	}

	return resource
}

func dbmateMigrate(pg *pgxpool.Pool) error {
	u, _ := url.Parse(connectionString)
	db := dbmate.New(u)
	cwd, _ := filepath.Abs(".")
	parentDir := filepath.Dir(cwd)
	migrationsDir := parentDir + "/db/migrations"

	db.MigrationsDir = migrationsDir

	if err := db.Migrate(); err != nil {
		return err
	}

	return nil
}

func seedDatabase(pg *pgxpool.Pool) error {
	if _, err := pg.Exec(context.Background(), `
		INSERT INTO users (name, lastname) VALUES ('andreas', 'hansson')
	`); err != nil {
		return err
	}

	return nil
}

func stopContainer(r *dockertest.Resource) {
	if err := r.Close(); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestMain(m *testing.M) {
	r := startPostgresPool()
	if err := dbmateMigrate(pg); err != nil {
		panic(err)
	}

	if err := seedDatabase(pg); err != nil {
		panic(err)
	}
	code := m.Run()
	stopContainer(r)
	tearDown()
	os.Exit(code)
}

func tearDown() {
	if err := os.Remove("db/schema.sql"); err != nil {
		panic(err)
	}
	if err := os.Remove("db"); err != nil {
		panic(err)
	}
}
