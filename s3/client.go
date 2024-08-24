package s3

import (
	"context"
	"fmt"
	cfg "goes-api-go/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	Client     *minio.Client
	BucketName string
	BaseURL    string
}

func NewS3Client(ctx context.Context, cfg *cfg.Config) (*S3Client, error) {
	// Create the MinIO client with the provided endpoint, access key, and secret key
	client, err := minio.New(cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
		return nil, err
	}

	// Construct the base URL for the S3 bucket
	baseURL := fmt.Sprintf("%s/%s/", cfg.S3Endpoint, cfg.BucketName)

	// Return the initialized S3Client
	return &S3Client{
		Client:     client,
		BucketName: cfg.BucketName,
		BaseURL:    baseURL,
	}, nil
}
