package awss3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	"github.com/google/uuid"
)

func UploadFileToS3(conf *config.Config, file multipart.File, fileName string) (string, string, error) {
	uniqueFileName := uuid.New().String() + filepath.Ext(fileName)

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &conf.AWSS3.BucketName,
		Key:    &uniqueFileName,
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead, // Set to PublicRead for public access
	})
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", conf.AWSS3.BucketName, uniqueFileName), uniqueFileName, nil
}

func GetPresignedURL(config *config.Config, fileName string) (string, error) {

	presignClient := s3.NewPresignClient(s3Client)

	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &config.AWSS3.BucketName,
		Key:    &fileName,
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", err
	}

	return req.URL, nil
}
