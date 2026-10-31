package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrNotReady = errors.New("storage: minio client not initialized")
var client *minio.Client

func envRequire(keys ...string) (string, error) {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}
	}
	return "", fmt.Errorf("storage: required env missing: %v", keys)
}

func Init(ctx context.Context) error {
	endpoint, err := envRequire("MINIO_ENDPOINT", "ANGEL_CACHE_YTI_SERVICE_SERVICE_HOST")
	if err != nil { return err }
	port, err := envRequire("MINIO_PORT", "ANGEL_CACHE_YTI_SERVICE_SERVICE_PORT_API")
	if err != nil { return err }
	ak, err := envRequire("MINIO_ROOT_USER", "MINIO_ACCESS_KEY")
	if err != nil { return err }
	sk, err := envRequire("MINIO_ROOT_PASSWORD", "MINIO_SECRET_KEY")
	if err != nil { return err }

	c, err := minio.New(endpoint+":"+port, &minio.Options{
		Creds:  credentials.NewStaticV4(ak, sk, ""),
		Secure: os.Getenv("MINIO_USE_TLS") == "1",
	})
	if err != nil { return fmt.Errorf("storage: minio.New: %w", err) }
	
	if _, err := c.ListBuckets(ctx); err != nil {
		return fmt.Errorf("storage: ListBuckets %s: %w", endpoint, err)
	}
	client = c
	log.Printf("[minio] OK endpoint=%s", endpoint)
	return nil
}

func Client() *minio.Client { return client }

func ClientOrErr() (*minio.Client, error) {
	if client == nil { return nil, ErrNotReady }
	return client, nil
}
