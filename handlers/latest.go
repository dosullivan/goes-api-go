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
	return func(c *gin.Context) {
		maxDaysBack := 7
		ctx := context.Background()
		var latestObjectKey string

		for i := 0; i < maxDaysBack; i++ {
			checkDate := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
			prefix := "false-color/fd/" + checkDate + "/"

			// Use Minio's ListObjects to list objects in the bucket with the specified prefix
			objectCh := s3Client.Client.ListObjects(ctx, s3Client.BucketName, minio.ListObjectsOptions{
				Prefix:    prefix,
				Recursive: true,
			})

			var lastObject minio.ObjectInfo
			found := false
			for object := range objectCh {
				if object.Err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
					return
				}
				lastObject = object
				found = true
			}

			if found {
				latestObjectKey = lastObject.Key
				break
			}
		}

		if latestObjectKey == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No recent images found"})
			return
		}

		imageURL := s3Client.BaseURL + latestObjectKey
		c.JSON(http.StatusOK, gin.H{"imageUrl": imageURL})
	}
}
