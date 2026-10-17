package profile

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/angel/angel/internal/storage"
	"github.com/minio/minio-go/v7"
)

func handleUpload(w http.ResponseWriter, r *http.Request) {
	uid, err := authFromReq(r)
	if err != nil { 
		log.Printf("[UPLOAD] Unauthorized attempt: %v", err)
		writeJSON(w, 401, map[string]any{"ok": false, "error": "unauthorized"}); 
		return 
	}
	
	log.Printf("[UPLOAD] Received upload request from user %d", uid)

	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil { 
		log.Printf("[UPLOAD] ParseMultipartForm failed: %v", err)
		writeJSON(w, 400, map[string]any{"ok": false, "error": "file too large"}); 
		return 
	}

	file, header, err := r.FormFile("image")
	if err != nil { 
		log.Printf("[UPLOAD] FormFile failed: %v", err)
		writeJSON(w, 400, map[string]any{"ok": false, "error": "missing image"}); 
		return 
	}
	defer file.Close()

	log.Printf("[UPLOAD] User %d uploading file: %s, Size: %d bytes", uid, header.Filename, header.Size)

	ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, "."):])
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		log.Printf("[UPLOAD] Invalid format: %s", ext)
		writeJSON(w, 400, map[string]any{"ok": false, "error": "invalid format"}); 
		return
	}

	objectName := fmt.Sprintf("users/%d_%d%s", uid, time.Now().Unix(), ext)
	bucketName := "avatars"

	log.Printf("[UPLOAD] Checking bucket %s existence...", bucketName)
	exists, err := storage.Client().BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Printf("[UPLOAD] BucketExists error: %v", err)
	}
	
	if !exists {
		log.Printf("[UPLOAD] Bucket %s does not exist, creating...", bucketName)
		err = storage.Client().MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil { 
			log.Printf("[UPLOAD] MakeBucket failed: %v", err)
			writeJSON(w, 500, map[string]any{"ok": false, "error": "bucket creation failed"}); 
			return 
		}
	}
	
	// ALWAYS set public policy (this fixes the image not showing bug)
	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`
	err = storage.Client().SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
	    log.Printf("[UPLOAD] SetBucketPolicy warning: %v", err)
	}

	log.Printf("[UPLOAD] Uploading %s to MinIO...", objectName)
	_, err = storage.Client().PutObject(context.Background(), bucketName, objectName, file, header.Size, minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")})
	if err != nil { 
		log.Printf("[UPLOAD] PutObject failed: %v", err)
		writeJSON(w, 500, map[string]any{"ok": false, "error": "upload failed: " + err.Error()}); 
		return 
	}

	fileURL := fmt.Sprintf("https://files.nomi98.com/%s/%s", bucketName, objectName)
	log.Printf("[UPLOAD] Success! User %d -> %s", uid, fileURL)
	writeJSON(w, 200, map[string]any{"ok": true, "url": fileURL})
}