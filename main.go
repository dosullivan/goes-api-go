package main

import (
	"context"
	"log"
	"strings"

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

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Check if TRUSTED_PROXIES is set or not
	if cfg.TrustedProxies == "" {
		// Disable proxy support by setting nil
		if err := router.SetTrustedProxies(nil); err != nil {
			log.Fatalf("Failed to set trusted proxies: %v", err)
		}
	} else {
		// Split the string into a slice of trusted proxy IP ranges
		trustedProxies := strings.Split(cfg.TrustedProxies, ",")
		log.Printf("Trusted proxies: %v", trustedProxies)
		if err := router.SetTrustedProxies(trustedProxies); err != nil {
			log.Fatalf("Failed to set trusted proxies: %v", err)
		}
	}

	router.GET("/latest", handlers.GetLatestImage(s3Client))
	router.GET("/archive/:date", handlers.GetImagesByDate(s3Client))
	router.GET("/available-dates", handlers.GetAvailableDates(s3Client))

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
