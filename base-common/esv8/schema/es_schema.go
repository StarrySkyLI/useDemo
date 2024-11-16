package schema

type EsSchema struct {
	id string // 主键key
}

func (s *EsSchema) GetId() string {
	return s.id
}

func (s *EsSchema) SetId(id string) {
	s.id = id
}
