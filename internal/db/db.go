package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

//go:embed migrations/*.sql
var migrationsFS embed.FS

func runMigrations(ctx context.Context, p *pgxpool.Pool) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, n := range names {
		data, err := migrationsFS.ReadFile("migrations/" + n)
		if err != nil {
			return err
		}
		log.Printf("[db] migration: %s", n)
		if _, err := p.Exec(ctx, string(data)); err != nil {
			return fmt.Errorf("%s: %w", n, err)
		}
	}
	log.Printf("[db] migrations: %d applied", len(names))
	return nil
}


func Init(ctx context.Context) error {
	host := first(
		os.Getenv("PG_HOST"),
		"angel-db-brp-service")
	port := "5432"
	user := first(os.Getenv("PGUSER"), "postgres")
	pass := os.Getenv("PostgreSQL")
	name := first(os.Getenv("POSTGRES_DB"), "angel-dbhsl_db")

	log.Printf("[db] env: BOT_PGF=%q BRP_HOST=%q PGUSER=%q DB=%q",
		os.Getenv("ANGEL_BOT_PGF_SERVICE_SERVICE_HOST"),
		os.Getenv("ANGEL_DB_BRP_SERVICE_SERVICE_HOST"),
		os.Getenv("PGUSER"),
		os.Getenv("POSTGRES_DB"))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
	log.Printf("[db] connecting to %s:%s/%s user=%s passlen=%d", host, port, name, user, len(pass))

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil { return fmt.Errorf("parse: %w", err) }
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = 5 * time.Minute

	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil { return fmt.Errorf("pool: %w", err) }

	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := p.Ping(pctx); err != nil { p.Close(); return fmt.Errorf("ping: %w", err) }

	if err := runMigrations(ctx, p); err != nil {
		p.Close()
		return fmt.Errorf("migrate: %w", err)
	}

	pool = p
	log.Printf("[db] OK")
	return nil
}

func Pool() *pgxpool.Pool { return pool }

func Close() { if pool != nil { pool.Close() } }

func first(ss ...string) string {
	for _, s := range ss { if s != "" { return s } }
	return ""
}
