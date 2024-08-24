package handlers

import (
	"context"
	"net/http"

	"goes-api-go/s3"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func GetImagesByDate(s3Client *s3.S3Client) gin.HandlerFunc {
	// Return a Gin handler function that processes the HTTP request
	return func(c *gin.Context) {
		// Extract the 'date' parameter from the URL (Expected format: YYYY-MM-DD)
		date := c.Param("date") // Expected format: YYYY-MM-DD
		// Construct the prefix for the directory structure in the S3 bucket based on the date
		prefix := "false-color/fd/" + date + "/"
		// Create a background context for managing the lifecycle of the request to the S3 storage
		ctx := context.Background()

		// Use Minio's ListObjects function to list objects in the bucket with the specified prefix
		objectCh := s3Client.Client.ListObjects(ctx, s3Client.BucketName, minio.ListObjectsOptions{
			Prefix:    prefix,  // Filter objects by the specified prefix
			Recursive: true,    // Recursive listing to include all objects under the prefix
		})

		// Initialize a slice to store the URLs of the images found
		var imageUrls []string
		// Iterate over the objects returned by ListObjects
		for object := range objectCh {
			// Check if there was an error retrieving the object
			if object.Err != nil {
				// If there's an error, respond with a 500 Internal Server Error and the error message
				c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
				return
			}
			// Construct the full URL for each object and append it to the imageUrls slice
			imageUrls = append(imageUrls, s3Client.BaseURL+object.Key)
		}

		// If no images were found for the specified date, respond with a 404 Not Found
		if len(imageUrls) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No images found for the specified date"})
			return
		}

		// Respond to the HTTP request with a JSON object containing the list of image URLs
		c.JSON(http.StatusOK, gin.H{"imageUrls": imageUrls})
	}
}
