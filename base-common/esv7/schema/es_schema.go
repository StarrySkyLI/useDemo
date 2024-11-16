package schema

import (
	esv7 "gitlab.coolgame.world/go-template/base-common/esv7"
)

type EsSchema struct {
	id string // 主键key
}

func (s *EsSchema) GetId() string {
	return s.id
}

func (s *EsSchema) SetId(id string) {
	s.id = id
}

var _ esv7.Schema = (*EsSchema)(nil)
