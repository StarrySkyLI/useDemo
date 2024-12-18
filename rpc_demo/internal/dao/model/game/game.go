package game

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"rpc_demo/internal/dao"
	"rpc_demo/internal/dao/dto"
	"rpc_demo/internal/dao/schema"
)

type (
	GameModel struct {
		ctx   context.Context
		data  *dao.Dao
		log   logx.Logger
		model *schema.Game
	}
	IGameModel interface {
		GetModel() *schema.Game

		GameList(param dto.GameList) (models []schema.Game, total int64, err error)

		FindOne(Id int64) (info schema.Game, err error)
	}
)

func NewGameModel(ctx context.Context, data *dao.Dao, log logx.Logger) IGameModel {
	return &GameModel{ctx, data, log, &schema.Game{}}
}
func (model *GameModel) GetModel() *schema.Game {

	return model.model
}

func (model *GameModel) GameList(param dto.GameList) (models []schema.Game, total int64, err error) {
	var count int64
	if err = model.data.DB.WithContext(model.ctx).
		Model(model.model).
		Count(&count).
		Order("created_at desc").
		Offset(int((param.Page - 1) * param.PageSize)).
		Limit(int(param.PageSize)).
		Find(&models).Error; err != nil {
		return nil, count, err
	}
	return models, count, nil
}

func (model *GameModel) FindOne(Id int64) (info schema.Game, err error) {
	if err = model.data.DB.WithContext(model.ctx).
		Model(model.model).
		Where("id = ?", Id).First(&info).Error; err != nil {
		return info, errors.Wrap(err, "mysql: FindOne")
	}
	return info, nil
}
