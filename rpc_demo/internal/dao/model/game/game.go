package game

import (
	"context"
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
	//TODO implement me
	panic("implement me")
}

func (model *GameModel) FindOne(Id int64) (info schema.Game, err error) {
	//TODO implement me
	panic("implement me")
}
