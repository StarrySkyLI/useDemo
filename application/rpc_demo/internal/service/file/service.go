package file

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"useDemo/application/rpc_demo/internal/dao/dto"
	"useDemo/application/rpc_demo/internal/dao/model/file"
	"useDemo/application/rpc_demo/internal/dao/schema"
	"useDemo/application/rpc_demo/internal/svc"
	"useDemo/application/rpc_demo/rpc"
	"useDemo/application/rpc_demo/utils"
	"useDemo/base-common/snowflake"
)

type FileService struct {
	ctx       context.Context
	svc       *svc.ServiceContext
	fileModel *file.FileModel
}

func NewFileService(ctx context.Context, svc *svc.ServiceContext) *FileService {
	return &FileService{
		ctx:       ctx,
		svc:       svc,
		fileModel: file.NewFileModel(ctx, svc.Dao, logx.WithContext(ctx)),
	}
}
func (s *FileService) CheckFileExistedAndGetFile(in *rpc.FileContext) (int64, bool, error) {
	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("hash = ?", in.Hash).Where("is_deleted = ?", false).Where("uploaded = ?", true)
	})
	if err != nil {
		return 0, false, err
	}
	return f.ID, f.ID > 0, nil

}

func (s *FileService) PreSignPut(in *rpc.FileContext) (string, int64, error) {
	f := &schema.File{
		ID:         snowflake.GetSnowflakeId(),
		DomainName: in.Domain,
		BizName:    in.BizName,
		Hash:       in.Hash,
		FileSize:   in.Size,
		FileType:   in.FileType,
		Uploaded:   false,
		IsDeleted:  false,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
	err := s.fileModel.CreateFile(f)
	if err != nil {
		return "", 0, err
	}
	url, err := s.svc.Minio.PreSignPutUrl(context.Background(), in.Domain, utils.GetObjectName(in.BizName, f.ID), in.ExpireSeconds)
	if err != nil {
		return "", 0, err
	}
	return url, f.ID, nil
}

func (s *FileService) ReportUploaded(in *rpc.FileContext) error {
	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", in.FileId)
	})
	if err != nil {
		return err
	}
	hash, err := s.svc.Minio.GetObjectHash(context.Background(), in.Domain, utils.GetObjectName(in.BizName, f.ID))
	if err != nil {
		return err
	}
	if !utils.CheckHash(hash, f.Hash) && !strings.Contains(hash, "-") {
		return errors.New("failed to validate hash of uploaded file")
	}
	f.Uploaded = true
	return s.fileModel.UpdateFile(&f)

}

func (s *FileService) PreGetFile(in *rpc.FileContext) (string, error) {
	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", in.FileId).Where("uploaded = ?", true)
	})
	if err != nil {
		return "", err
	}
	if f.ID == 0 {
		return "", errors.New("file not found")
	}
	return s.svc.Minio.PreSignGetUrl(context.Background(), in.Domain, utils.GetObjectName(f.BizName, f.ID), in.Filename, in.ExpireSeconds)

}

func (s *FileService) PreSignGet(fileContext *rpc.FileContext) (string, error) {
	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", fileContext.FileId).Where("uploaded = ?", true)
	})
	if err != nil {
		return "", err
	}
	return s.svc.Minio.PreSignGetUrl(context.Background(), f.DomainName, utils.GetObjectName(f.BizName, f.ID), fileContext.Filename, fileContext.ExpireSeconds)
}

func (s *FileService) PreSignSlicingPut(in *rpc.FileContext) (*dto.SlicingFile, error) {
	f := &schema.File{
		ID:         snowflake.GetSnowflakeId(),
		DomainName: in.Domain,
		BizName:    in.BizName,
		Hash:       in.Hash,
		FileSize:   in.Size,
		FileType:   in.FileType,
		Uploaded:   false,
		IsDeleted:  false,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
	err := s.fileModel.CreateFile(f)
	if err != nil {
		return nil, err
	}
	uploadId, err := s.svc.Minio.CreateSlicingUpload(context.Background(), in.Domain, utils.GetObjectName(in.BizName, f.ID), minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	slicingFile := dto.New(f).SetUploadId(uploadId).SetTotalParts()
	urls := make([]string, slicingFile.TotalParts)
	for i := 1; i <= int(slicingFile.TotalParts); i++ {
		url, e := s.svc.Minio.PreSignSlicingPutUrl(context.Background(), f.DomainName, utils.GetObjectName(f.BizName, f.ID), uploadId, int64(i))
		if e != nil {
			return nil, e
		}
		urls[i-1] = url
	}
	slicingFile.UploadUrl = urls
	return slicingFile, nil
}

func (s *FileService) MergeFileParts(in *rpc.MergeFilePartsRequest) error {
	uploadResult, err := s.GetProgressRate4SlicingPut(in.UploadId, in.FileContext)
	if err != nil {
		return err
	}

	if ok, _ := s.checkSlicingFileUploaded(uploadResult); !ok {
		return errors.New("not all parts uploaded")
	}

	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", in.FileContext.FileId)
	})
	if err != nil {
		return err
	}
	sf := dto.New(&f).SetTotalParts()

	result, err := s.svc.Minio.ListSlicingFileParts(context.Background(), f.DomainName, utils.GetObjectName(f.BizName, f.ID), in.UploadId, sf.TotalParts)
	if err != nil {
		return err
	}

	parts := make([]minio.CompletePart, 0)
	for i := 0; i < len(result.ObjectParts); i++ {
		parts = append(parts, minio.CompletePart{
			PartNumber: i + 1,
			ETag:       result.ObjectParts[i].ETag,
		})
	}

	return s.svc.Minio.MergeSlices(context.Background(), f.DomainName, utils.GetObjectName(f.BizName, f.ID), in.UploadId, parts)
}

func (s *FileService) GetProgressRate4SlicingPut(uploadId string, fileContext *rpc.FileContext) (map[string]bool, error) {
	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", fileContext.FileId)
	})
	if err != nil {
		return nil, err
	}
	sf := dto.New(&f).SetTotalParts()
	result, err := s.svc.Minio.ListSlicingFileParts(context.Background(), f.DomainName, utils.GetObjectName(f.BizName, f.ID), uploadId, sf.TotalParts)
	if err != nil {
		return nil, err
	}

	res := make(map[string]bool)
	parts := result.ObjectParts
	for i := 0; i < int(sf.TotalParts); i++ {
		if len(parts[i].ETag) > 0 {
			res[strconv.FormatInt(int64(i+1), 10)] = true
		} else {
			res[strconv.FormatInt(int64(i+1), 10)] = false
		}
	}

	return res, nil
}

func (s *FileService) checkSlicingFileUploaded(result map[string]bool) (bool, string) {
	total := 0
	finished := 0
	for _, uploaded := range result {
		if uploaded {
			finished++
		}

		total++
	}

	rate := fmt.Sprintf("%d/%d", finished, total)
	return total == finished, rate
}

func (s *FileService) GetInfoById(domainName, bizName string, fileId int64) (*schema.File, error) {
	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", fileId)
	})
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *FileService) RemoveFile(fileContext *rpc.FileContext) error {
	f, err := s.fileModel.FirstByWhere(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", fileContext.FileId)
	})
	if err != nil {
		return err
	}
	if f.ID == 0 {
		return errors.New("file not found")
	}
	f.IsDeleted = true
	return s.fileModel.UpdateFile(&f)
}
