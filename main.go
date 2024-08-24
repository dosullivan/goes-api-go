package main

import (
	"context"
	"log"

	"goes-api-go/config"
	"goes-api-go/handlers"
	"goes-api-go/s3"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
)

func main() {
    ctx := context.Background()

    var cfg config.Config
    if err := envconfig.Process(ctx, &cfg); err != nil {
        log.Fatalf("Failed to load environment configuration: %v", err)
    }

    s3Client, err := s3.NewS3Client(ctx, &cfg)
    if err != nil {
        log.Fatalf("Failed to initialize S3 client: %v", err)
    }

    router := gin.Default()

    router.GET("/latest", handlers.GetLatestImage(s3Client))
    router.GET("/archive/:date", handlers.GetImagesByDate(s3Client))
    router.GET("/available-dates", handlers.GetAvailableDates(s3Client))

    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
