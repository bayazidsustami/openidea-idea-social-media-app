package security

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

func GetAws3Session() *s3.S3 {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("S3_REGION"),
		Credentials: credentials.NewStaticCredentials(
			viper.GetString("S3_ID"),
			viper.GetString("S3_SECRET_KEY"),
			"",
		),
	}))
	svc := s3.New(sess)

	return svc
}
