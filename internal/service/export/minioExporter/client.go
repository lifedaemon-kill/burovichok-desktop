package minioExporter

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

type Client struct {
	client     *minio.Client
	bucketName string
	zLog       logger.Logger
}

func NewClient(ctx context.Context, conf config.MinioConf, zLog logger.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	client, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		zLog.Errorw("minio.New failed", "error", err)
		return nil, errors.Wrap(err, "failed to initialize minio client")
	}

	bucketName := conf.BucketName
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		zLog.Errorw("minio.BucketExists failed", "bucket", bucketName, "error", err)
		return nil, errors.Wrapf(err, "failed to check if minio bucket '%s' exists", bucketName)
	}
	if !exists {
		zLog.Infow("MinIO bucket does not exist, attempting to create", "bucket", bucketName)
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}) // Можно указать регион, если нужно
		if err != nil {
			zLog.Errorw("minio.MakeBucket failed", "bucket", bucketName, "error", err)
			return nil, errors.Wrapf(err, "failed to create minio bucket '%s'", bucketName)
		}
		zLog.Infow("MinIO bucket created successfully", "bucket", bucketName)
	} else {
		zLog.Infow("MinIO bucket already exists", "bucket", bucketName)
	}
	zLog.Infow("MinIO client initialized successfully")

	return &Client{
		client:     client,
		bucketName: conf.BucketName,
		zLog:       zLog,
	}, nil
}

func (c *Client) Upload(ctx context.Context, basename string, buf *bytes.Buffer) (minio.UploadInfo, error) {
	// 4. Генерируем имя файла для MinIO
	timestamp := time.Now().Format("20060102150405") // Формат YYYYMMDD_HHMMSS

	objectName := fmt.Sprintf("%s_%s.zip", basename, timestamp)
	// 5. Загружаем архив в MinIO
	contentType := "application/zip"
	uploadInfo, err := c.client.PutObject(
		ctx,
		c.bucketName,
		objectName,       // Имя объекта в MinIO
		buf,              // Данные из буфера
		int64(buf.Len()), // Размер данных
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		c.zLog.Errorw("Failed to upload archive to MinIO", "object", objectName, "error", err)
		return minio.UploadInfo{}, errors.Wrapf(err, "failed to upload '%s' to bucket '%s'", objectName, c.bucketName)
	}

	c.zLog.Infow("Successfully uploaded archive to MinIO",
		"object", objectName,
		"bucket", uploadInfo.Bucket,
		"etag", uploadInfo.ETag,
		"size", uploadInfo.Size,
	)
	return uploadInfo, nil
}
