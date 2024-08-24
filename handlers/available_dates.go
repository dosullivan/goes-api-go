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
	return func(c *gin.Context) {
		prefix := "false-color/fd/"
		ctx := context.Background()

		// Use Minio's ListObjects to list objects in the bucket with the specified prefix
		objectCh := s3Client.Client.ListObjects(ctx, s3Client.BucketName, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: false, // Use recursive false to simulate delimiter behavior
		})

		datesSet := make(map[string]struct{})
		for object := range objectCh {
			if object.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
				return
			}

			parts := strings.Split(object.Key, "/")
			if len(parts) >= 4 {
				date := parts[2]
				datesSet[date] = struct{}{}
			}
		}

		dates := make([]string, 0, len(datesSet))
		for date := range datesSet {
			dates = append(dates, date)
		}

		c.JSON(http.StatusOK, gin.H{"availableDates": dates})
	}
}
