package handlers

import (
	"context"
	"net/http"
	"strings"

	"goes-api-go/s3"

	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
)

func GetAvailableDates(s3Client *s3.S3Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        prefix := "false-color/fd/"
        ctx := context.Background()

        resp, err := s3Client.Client.ListObjectsV2(ctx, &awsS3.ListObjectsV2Input{
            Bucket:    &s3Client.BucketName,
            Prefix:    &prefix,
            Delimiter: aws.String("/"),
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        datesSet := make(map[string]struct{})
        for _, commonPrefix := range resp.CommonPrefixes {
            parts := strings.Split(*commonPrefix.Prefix, "/")
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
