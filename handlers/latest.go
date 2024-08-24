package handlers

import (
	"context"
	"net/http"
	"time"

	"goes-api-go/s3"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func GetLatestImage(s3Client *s3.S3Client) gin.HandlerFunc {
	// Return a Gin handler function that processes the HTTP request
	return func(c *gin.Context) {
		maxDaysBack := 7            // Number of days to look back for the latest image
		ctx := context.Background() // Create a background context for managing the request lifecycle
		var latestObjectKey string  // Variable to store the key of the latest object found

		// Iterate over the last `maxDaysBack` days to find the most recent image
		for i := 0; i < maxDaysBack; i++ {
			// Calculate the date for each iteration (going back one day at a time)
			checkDate := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			// Construct the prefix for the directory structure in the S3 bucket
			prefix := "false-color/fd/" + checkDate + "/"

			// Use Minio's ListObjects to list objects in the bucket with the specified prefix
			objectCh := s3Client.Client.ListObjects(ctx, s3Client.BucketName, minio.ListObjectsOptions{
				Prefix:    prefix,  // Filter objects by the specified prefix
				Recursive: true,    // Recursive listing to include all objects under the prefix
			})

			var lastObject minio.ObjectInfo // Variable to hold the last object found
			found := false                  // Flag to indicate if any objects were found

			// Iterate over the objects returned by ListObjects
			for object := range objectCh {
                // Check if there was an error retrieving the object
				if object.Err != nil {
					// If there's an error, respond with a 500 Internal Server Error and the error message
					c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
					return
				}
				// Set the lastObject to the current object
				lastObject = object
				found = true  // Set the flag to true indicating an object was found
			}

			// If we found at least one object, store its key and stop searching further
			if found {
				latestObjectKey = lastObject.Key
				break
			}
		}

		// If no objects were found across the checked dates, respond with a 404 Not Found
		if latestObjectKey == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No recent images found"})
			return
		}

		// Construct the full URL for the latest image found
		imageURL := s3Client.BaseURL + latestObjectKey
		// Respond to the HTTP request with a JSON object containing the image URL
		c.JSON(http.StatusOK, gin.H{"imageUrl": imageURL})
	}
}
