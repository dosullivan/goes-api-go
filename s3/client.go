package s3

import (
	"context"
	"fmt"
	cfg "goes-api-go/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client     *s3.Client
	BucketName string
	BaseURL    string
}

func NewS3Client(ctx context.Context, cfg *cfg.Config) (*S3Client, error) {
	// Create the credentials provider using the credentials from the config
	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	// Create the custom endpoint resolver using the S3Endpoint from the config
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL: cfg.S3Endpoint,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested for service %s in region %s", service, region)
	})

	// Load the AWS configuration with the custom resolver and credentials provider
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create the S3 client using the custom endpoint resolver
	client := s3.NewFromConfig(awsConfig)

	// Construct the base URL for the S3 bucket
	baseURL := fmt.Sprintf("%s/%s/", cfg.S3Endpoint, cfg.BucketName)

	// Return the initialized S3Client
	return &S3Client{
		Client:     client,
		BucketName: cfg.BucketName,
		BaseURL:    baseURL,
	}, nil
}
