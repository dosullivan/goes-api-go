package handlers

import (
	"context"
	"net/http"

	"goes-api-go/s3"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func GetImagesByDate(s3Client *s3.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		date := c.Param("date") // Expected format: YYYY-MM-DD
		prefix := "false-color/fd/" + date + "/"
		ctx := context.Background()

		// Use Minio's ListObjects function to list objects in the bucket with the specified prefix
		objectCh := s3Client.Client.ListObjects(ctx, s3Client.BucketName, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		})

		var imageUrls []string
		for object := range objectCh {
			if object.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
				return
			}
			imageUrls = append(imageUrls, s3Client.BaseURL+object.Key)
		}

		if len(imageUrls) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No images found for the specified date"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"imageUrls": imageUrls})
	}
}
