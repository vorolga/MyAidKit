package composites

import (
	"context"
	"main/internal/constants"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioComposite struct {
	client *minio.Client
}

func NewMinioComposite() (*MinioComposite, error) {
	minioClient, err := minio.New(os.Getenv("MINIOURL"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIOUSER"), os.Getenv("MINIOPASSWORD"), ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(context.Background(), constants.UserObjectsBucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), constants.UserObjectsBucketName, minio.MakeBucketOptions{
			ObjectLocking: false,
		})
		if err != nil {
			return nil, err
		}
	}

	imageName := constants.DefaultImage
	file, err := os.Open(imageName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	opts := minio.PutObjectOptions{
		ContentType:  "image/png",
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	_, err = minioClient.PutObject(
		context.Background(),
		constants.UserObjectsBucketName, // Константа с именем бакета
		imageName,
		file,
		stat.Size(),
		opts,
	)

	if err != nil {
		return nil, err
	}

	exists, err = minioClient.BucketExists(context.Background(), constants.MedicinesObjectsBucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), constants.MedicinesObjectsBucketName, minio.MakeBucketOptions{
			ObjectLocking: false,
		})
		if err != nil {
			return nil, err
		}
	}

	imageName = constants.DefaultMedicine
	file, err = os.Open(imageName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err = file.Stat()
	if err != nil {
		return nil, err
	}

	opts = minio.PutObjectOptions{
		ContentType:  "image/png",
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	_, err = minioClient.PutObject(
		context.Background(),
		constants.MedicinesObjectsBucketName, // Константа с именем бакета
		imageName,
		file,
		stat.Size(),
		opts,
	)

	if err != nil {
		return nil, err
	}

	return &MinioComposite{client: minioClient}, nil
}
