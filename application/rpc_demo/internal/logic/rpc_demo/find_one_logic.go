package rpcdemologic

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"useDemo/application/rpc_demo/client/rpc_demo"
	"useDemo/application/rpc_demo/internal/dao/model/game"
	"useDemo/application/rpc_demo/rpc"

	"useDemo/application/rpc_demo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	gameModel game.IGameModel
}

var count int64

func NewFindOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneLogic {
	return &FindOneLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		gameModel: game.NewGameModel(ctx, svcCtx.Dao, logx.WithContext(ctx)),
	}
}

const (
	prefixGames = "biz#games#%d#%d"
)

func (l *FindOneLogic) FindOne(in *rpc_demo.GameInfoReq) (*rpc_demo.GameInfoRep, error) {
	// 熔断器测试
	atomic.AddInt64(&count, 1)
	if atomic.LoadInt64(&count)%2 == 0 {
		time.Sleep(time.Second)
		return nil, fmt.Errorf("timeout")
	}
	opProducts, ok := l.svcCtx.LocalCache.Get(gameKey(in.GetId(), 1))
	if ok {
		fmt.Println(opProducts)
		return &rpc_demo.GameInfoRep{
			Info: opProducts.(*rpc_demo.GameInfo),
		}, nil
	}
	info, err := l.gameModel.FindOne(in.GetId())
	if err != nil {
		return nil, err
	}
	resp := &rpc.GameInfo{
		Id:         info.ID,
		Name:       info.Name,
		Code:       info.Code,
		Data:       info.Data,
		CreateTime: strconv.FormatInt(info.CreatedAt.Unix(), 10),
		UpdateTime: strconv.FormatInt(info.UpdatedAt.Unix(), 10),
	}

	l.svcCtx.LocalCache.Set(gameKey(in.GetId(), 1), resp)

	return &rpc_demo.GameInfoRep{
		Info: resp,
	}, nil
}
func gameKey(uid int64, sortType int32) string {
	return fmt.Sprintf(prefixGames, uid, sortType)
}
