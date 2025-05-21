package minio_service

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioConf struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}
type MinioService struct {
	core *minio.Core
}

func NewMinio(core *minio.Core) *MinioService {
	return &MinioService{
		core: core,
	}
}

func (r *MinioService) PreSignGetUrl(ctx context.Context, bucketName, objectName, fileName string, expireSeconds int64) (string, error) {
	reqParams := make(url.Values)

	if fileName != "" {
		reqParams.Set("response-content-disposition", "attachment; filename="+fileName)
	}

	signedUrl, err := r.core.PresignedGetObject(ctx, bucketName, objectName, time.Duration(expireSeconds*1000000000), reqParams)
	if err != nil {
		return "", err
	}

	return signedUrl.String(), nil
}

func (r *MinioService) PreSignPutUrl(ctx context.Context, bucketName, objectName string, expireSeconds int64) (string, error) {
	signedUrl, err := r.core.PresignedPutObject(ctx, bucketName, objectName, time.Duration(expireSeconds*1000000000))
	if err != nil {
		return "", err
	}

	return signedUrl.String(), nil
}

func (r *MinioService) CreateSlicingUpload(ctx context.Context, bucketName, objectName string, options minio.PutObjectOptions) (uploadId string, err error) {
	return r.core.NewMultipartUpload(ctx, bucketName, objectName, options)
}

func (r *MinioService) ListSlicingFileParts(ctx context.Context, bucketName, objectName, uploadId string, partsNum int64) (minio.ListObjectPartsResult, error) {
	var nextPartNumberMarker int
	return r.core.ListObjectParts(ctx, bucketName, objectName, uploadId, nextPartNumberMarker, int(partsNum)+1)
}

func (r *MinioService) PreSignSlicingPutUrl(ctx context.Context, bucketName, objectName, uploadId string, parts int64) (string, error) {
	params := url.Values{
		"uploadId":   {uploadId},
		"partNumber": {strconv.FormatInt(parts, 10)},
	}

	signedUrl, err := r.core.Presign(ctx, http.MethodPut, bucketName, objectName, time.Hour, params)
	if err != nil {
		return "", err
	}

	return signedUrl.String(), nil
}

func (r *MinioService) MergeSlices(ctx context.Context, bucketName, objectName, uploadId string, parts []minio.CompletePart) error {
	_, err := r.core.CompleteMultipartUpload(ctx, bucketName, objectName, uploadId, parts, minio.PutObjectOptions{})
	return err
}

func (r *MinioService) GetObjectHash(ctx context.Context, bucketName, objectName string) (string, error) {
	stats, err := r.core.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return "", err
	}

	return strings.ToUpper(stats.ETag), nil
}
