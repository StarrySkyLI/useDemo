package file

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"useDemo/application/rpc_demo/internal/dao"
	"useDemo/application/rpc_demo/internal/dao/schema"
)

type (
	FileModel struct {
		ctx   context.Context
		data  *dao.Dao
		log   logx.Logger
		model *schema.File
	}
)

func NewFileModel(ctx context.Context, data *dao.Dao, log logx.Logger) *FileModel {
	return &FileModel{ctx, data, log, &schema.File{}}
}
func (m FileModel) CreateFile(file *schema.File) error {
	err := m.data.DB.Model(m.model).Create(file).Error
	return err
}

func (m FileModel) FirstByWhere(any func(tx *gorm.DB) *gorm.DB) (schema.File, error) {
	var file schema.File
	dbRes := m.data.DB.Model(m.model)
	dbRes = any(dbRes)
	err := dbRes.First(&file).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		m.log.Error("mysql: FindOne :", err)
		return file, err
	}
	return file, nil
}

func (m FileModel) UpdateFile(f *schema.File) error {
	return m.data.DB.Model(m.model).Where("id = ?", f.ID).Updates(f).Error
}
