package filesystem

import (
	"context"
	"fmt"
	"go-upcycle_connect-backend/internal"
	"go-upcycle_connect-backend/utils/log"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UploadParams struct {
	Key     string
	Content string
}

func ListRustFS(client *s3.Client) []string {
	ctx := context.TODO()
	resp, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}

	buckets := make([]string, len(resp.Buckets))
	for i, b := range resp.Buckets {
		buckets[i] = *b.Name
	}
	return buckets
}

func ListObjects(bucket internal.RustFSBucket) {
	ctx := context.TODO()
	resp, err := bucket.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket.BucketName),
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, obj := range resp.Contents {
		fmt.Println(" -", *obj.Key)
	}
}

func UploadObject(bucket internal.RustFSBucket, params UploadParams) {
	ctx := context.TODO()
	_, err := bucket.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(params.Key),
		Body:   strings.NewReader(params.Content),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(fmt.Sprintf("File uploaded to %s/%s", bucket.BucketName, params.Key))
}

func GetObject(bucket internal.RustFSBucket, id string) (string, error) {
	ctx := context.TODO()
	resp, err := bucket.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return "", fmt.Errorf("download object failed: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read object content failed: %v", err)
	}
	return string(data), nil
}

func DeleteObject(bucket internal.RustFSBucket, key string) error {
	ctx := context.TODO()
	_, err := bucket.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete object failed: %v", err)
	}
	log.Info(fmt.Sprintf("File deleted from %s/%s", bucket.BucketName, key))
	return nil
}
