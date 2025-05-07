package game

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"useDemo/application/rpc_demo/internal/dao"
	"useDemo/application/rpc_demo/internal/dao/dto"
	"useDemo/application/rpc_demo/internal/dao/schema"
)

const (
	prefixGames = "biz#games#%d#%d"
	gamesExpire = 3600 * 24 * 2
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

		UpdateOne(info schema.Game) (err error)
	}
)

func (model *GameModel) UpdateOne(info schema.Game) (err error) {
	// TODO implement me
	panic("implement me")
}

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
	var isCache bool

	game, _ := model.cacheGame(model.ctx, Id)
	if len(game) > 0 {
		isCache = true
		err := json.Unmarshal([]byte(game), &info)
		if err != nil {
			logc.Errorf(model.ctx, "json.Unmarshal error: %v", err)
		}
	} else {
		if err = model.data.DB.WithContext(model.ctx).
			Model(model.model).
			Where("id = ?", Id).First(&info).Error; err != nil {
			return info, errors.Wrap(err, "mysql: FindOne")
		}
	}
	if !isCache {
		threading.GoSafe(func() {
			err = model.addGames(context.Background(), info, Id)
			if err != nil {
				logc.Errorf(model.ctx, "addCacheArticles error: %v", err)
			}
		})
	}
	return info, nil
}
func (model *GameModel) cacheGame(ctx context.Context, uid int64) (string, error) {
	key := gameKey(uid, 1)
	existsCtx, _ := model.data.BizRedis.ExistsCtx(ctx, key)
	if existsCtx {
		// 防止缓存击穿,缓存击穿经常发生在热点数据过期失效的时候，那么我们不让缓存失效不就好了，每次查询缓存
		// 的时候使用Exists来判断key是否存在，如果存在就使用Expire给缓存续期，既然是热点数据通过不断
		// 地续期也就不会过期了
		err := model.data.BizRedis.ExpireCtx(ctx, key, gamesExpire)
		if err != nil {
			logc.Error(ctx, "ExpireCtx key: %s error: %v", uid, err)
		}
	}
	res, err := model.data.BizRedis.GetCtx(model.ctx, key)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (model *GameModel) addGames(ctx context.Context, info schema.Game, id int64) error {
	key := gameKey(id, 1)
	value, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = model.data.BizRedis.SetCtx(ctx, key, string(value))
	if err != nil {
		return err
	}
	return model.data.BizRedis.ExpireCtx(ctx, key, gamesExpire)
}
func gameKey(uid int64, sortType int32) string {
	return fmt.Sprintf(prefixGames, uid, sortType)
}
