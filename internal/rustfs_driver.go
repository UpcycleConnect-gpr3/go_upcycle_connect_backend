package internal

import (
	"context"
	"fmt"
	"go-upcycle_connect-backend/utils/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type RustFSBucket struct {
	Client     *s3.Client
	BucketName string
}

func NewRustFS(
	bucketName, region, accessKeyID, secretAccessKey, endpoint string,
) RustFSBucket {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(endpoint)
	})

	log.Info(fmt.Sprintf("(CONFIG) RustFS Driver Initialized %s/%s", endpoint, bucketName))
	return RustFSBucket{Client: client, BucketName: bucketName}
}
