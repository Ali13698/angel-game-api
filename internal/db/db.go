package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Init(ctx context.Context) error {
	host := first(
		os.Getenv("ANGEL_BOT_PGF_SERVICE_SERVICE_HOST"),
		os.Getenv("ANGEL_DB_BRP_SERVICE_SERVICE_HOST"),
		"10.98.159.199")
	port := first(
		os.Getenv("ANGEL_BOT_PGF_SERVICE_SERVICE_PORT"),
		os.Getenv("ANGEL_DB_BRP_SERVICE_SERVICE_PORT"),
		"5432")
	user := first(os.Getenv("PGUSER"), "postgres")
	pass := first(os.Getenv("PostgreSQL2"), os.Getenv("PostgreSQL"), "")
	name := first(os.Getenv("POSTGRES_DB"), "angel-dbhs1_db")

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
