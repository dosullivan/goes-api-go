package handlers

import (
	"context"
	"net/http"

	"goes-api-go/s3"

	"github.com/gin-gonic/gin"

	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetImagesByDate(s3Client *s3.S3Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        date := c.Param("date") // Expected format: YYYY-MM-DD
        prefix := "false-color/fd/" + date + "/"
        ctx := context.Background()

        resp, err := s3Client.Client.ListObjectsV2(ctx, &awsS3.ListObjectsV2Input{
            Bucket: &s3Client.BucketName,
            Prefix: &prefix,
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if len(resp.Contents) == 0 {
            c.JSON(http.StatusNotFound, gin.H{"message": "No images found for the specified date"})
            return
        }

        imageUrls := make([]string, len(resp.Contents))
        for i, item := range resp.Contents {
            imageUrls[i] = s3Client.BaseURL + *item.Key
        }

        c.JSON(http.StatusOK, gin.H{"imageUrls": imageUrls})
    }
}
