package esv8

type IndexTable interface {
	IndexName() string
}

type Schema interface {
	GetId() string
	SetId(id string)
}
