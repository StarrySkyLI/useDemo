package dbM

import (
	"context"
	"fmt"

	"gitlab.coolgame.world/go-template/base-common/rpc"
	"gorm.io/gorm"
)

type SelectDbConfig struct {
	BusinessCode string
	DataSource   string
}

type SelectDb struct {
	dbMap map[string]*gorm.DB
}

func NewSelectDb(config []SelectDbConfig) *SelectDb {
	res := &SelectDb{
		dbMap: make(map[string]*gorm.DB),
	}

	for _, v := range config {
		res.dbMap[v.BusinessCode] = DbConnect(v.DataSource)
	}

	return res
}

func (s *SelectDb) GetDb(ctx context.Context) *gorm.DB {
	key := rpc.GetBusinessCode(ctx)
	if key == "" {
		panic("Business Code Not Exists. ")
	}

	if db, ok := s.dbMap[key]; ok {
		return db
	}

	panic(fmt.Sprintf("Business[%s] no DB", key))
}
