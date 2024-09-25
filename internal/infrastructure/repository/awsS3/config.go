package awss3

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
)

var s3Client *s3.Client

func InitS3(conf *cfg.Config) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(conf.AWSS3.Region))
	if err != nil {
		log.Fatal(err)
		return err
	}

	s3Client = s3.NewFromConfig(cfg)

	return nil
}
