package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

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

	router.GET("/proxy/image", func(c *gin.Context) {
		imageUrl := c.Query("url")
		if imageUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image URL is required"})
			return
		}

		// Validate URL is from our Minio instance
		if !strings.HasPrefix(imageUrl, s3Client.BaseURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image source"})
			return
		}

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(imageUrl)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
			return
		}
		defer resp.Body.Close()

		// Set appropriate headers for image delivery
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		c.Header("Cache-Control", "public, max-age=3600") // Optional: Add caching

		c.Stream(func(w io.Writer) bool {
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				log.Printf("Error while streaming image: %v", err)
				return false // Stop streaming on error
			}
			return false // Stop streaming after successful copy
		})
	})

	// Run the server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
