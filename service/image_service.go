package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type ImageService interface {
	UploadImage(ctx context.Context, file multipart.File, fileName string) (string, error)
}

type ImageServiceImpl struct {
	Svc *s3.S3
}

func NewImageService(
	svc *s3.S3,
) ImageService {
	return &ImageServiceImpl{
		Svc: svc,
	}
}

func (service *ImageServiceImpl) UploadImage(ctx context.Context, file multipart.File, fileName string) (string, error) {
	timeOut := 15 * time.Second

	var cancleFn func()
	if timeOut > 0 {
		ctx, cancleFn = context.WithTimeout(ctx, timeOut)
	}

	if cancleFn != nil {
		defer cancleFn()
	}

	contentType := http.DetectContentType([]byte(fileName))

	_, err := service.Svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(viper.GetString("S3_BUCKET_NAME")),
		Key:         aws.String(fileName),
		ACL:         aws.String("public-read"),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", viper.GetString("S3_BUCKET_NAME"), viper.GetString("S3_REGION"), fileName)
	return url, nil
}
