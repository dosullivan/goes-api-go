package handlers

import (
	"context"
	"net/http"
	"strings"

	"goes-api-go/s3"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func GetAvailableDates(s3Client *s3.S3Client) gin.HandlerFunc {
	// Return a Gin handler function that processes the HTTP request
	return func(c *gin.Context) {
		// Prefix for the directory structure within the S3 bucket where images are stored
		prefix := "false-color/fd/"
		// Create a background context for managing the lifecycle of the request to the S3 storage
		ctx := context.Background()

		// Use Minio's ListObjects function to list objects in the bucket under the specified prefix
		objectCh := s3Client.Client.ListObjects(ctx, s3Client.BucketName, minio.ListObjectsOptions{
			Prefix:    prefix,      // Filter objects by the specified prefix
			Recursive: false,       // Non-recursive listing to simulate AWS S3 delimiter behavior
		})

		// Create a set to store unique dates found in the object keys
		datesSet := make(map[string]struct{})
		// Iterate over the objects returned by ListObjects
		for object := range objectCh {
			if object.Err != nil {
				// If there's an error, respond with a 500 Internal Server Error and the error message
				c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
				return
			}

			// Split the object key by "/" to extract parts of the directory structure
			parts := strings.Split(object.Key, "/")
			// Ensure the key has the expected structure (at least 4 parts: "false-color", "fd", "YYYY-MM-DD", "")
			if len(parts) >= 4 {
				// Extract the date portion from the key (expected to be the third part)
				date := parts[2]
				// Store the date in the set to ensure uniqueness
				datesSet[date] = struct{}{}
			}
		}

		// Convert the set of unique dates into a slice for easier processing and JSON serialization
		dates := make([]string, 0, len(datesSet))
		for date := range datesSet {
			dates = append(dates, date)
		}

		// Respond to the HTTP request with a JSON object containing the list of available dates
		c.JSON(http.StatusOK, gin.H{"availableDates": dates})
	}
}
