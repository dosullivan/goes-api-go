package handlers

import (
	"context"
	"net/http"
	"time"

	"goes-api-go/s3"

	"github.com/gin-gonic/gin"

	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetLatestImage(s3Client *s3.S3Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        maxDaysBack := 7
        ctx := context.Background()
        var latestObjectKey string

        for i := 0; i < maxDaysBack; i++ {
            checkDate := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
            prefix := "false-color/fd/" + checkDate + "/"

            resp, err := s3Client.Client.ListObjectsV2(ctx, &awsS3.ListObjectsV2Input{
                Bucket: &s3Client.BucketName,
                Prefix: &prefix,
            })
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            if len(resp.Contents) > 0 {
                latestObject := resp.Contents[len(resp.Contents)-1]
                latestObjectKey = *latestObject.Key
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
