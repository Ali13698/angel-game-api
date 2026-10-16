package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

func Init(ctx context.Context) error {
	endpoint := first(os.Getenv("ANGEL_CACHE_YTI_SERVICE_SERVICE_HOST"), "127.0.0.1")
	port := first(os.Getenv("ANGEL_CACHE_YTI_SERVICE_SERVICE_PORT_API"), "9000")
	ak := first(os.Getenv("Minio"), os.Getenv("MINIO_ACCESS_KEY"), "admin")
	sk := first(os.Getenv("MINIO_ROOT_PASSWORD"), os.Getenv("MINIO_ROOT_USER"), "")
	c, err := minio.New(endpoint+":"+port, &minio.Options{
		Creds: credentials.NewStaticV4(ak, sk, ""), Secure: false,
	})
	if err != nil { return err }
	if _, err := c.ListBuckets(ctx); err != nil { return err }
	client = c
	log.Printf("[minio] OK")
	return nil
}

func Client() *minio.Client { return client }

func first(ss ...string) string {
	for _, s := range ss { if s != "" { return s } }
	return ""
}
