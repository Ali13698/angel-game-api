package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func Init(ctx context.Context) error {
	host := first(os.Getenv("ANGEL_REDIS_BUJ_SERVICE_SERVICE_HOST"), "127.0.0.1")
	port := first(os.Getenv("ANGEL_REDIS_BUJ_SERVICE_SERVICE_PORT"), "6379")
	pass := os.Getenv("REDIS")
	rdb = redis.NewClient(&redis.Options{Addr: host + ":" + port, Password: pass, DB: 0})
	if err := rdb.Ping(ctx).Err(); err != nil { return err }
	log.Printf("[redis] OK")
	return nil
}

func Client() *redis.Client { return rdb }
func Close() { if rdb != nil { _ = rdb.Close() } }

func first(ss ...string) string {
	for _, s := range ss { if s != "" { return s } }
	return ""
}
