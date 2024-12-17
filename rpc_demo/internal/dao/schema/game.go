package schema

import "time"

const TableNameGame = "game"

type Game struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键ID" json:"id"` // 主键ID

	Name      string    `gorm:"column:name;not null;comment:操作类型" json:"name"`
	Code      string    `gorm:"column:code;not null;comment:用户类型" json:"code"`
	Data      string    `gorm:"column:data;not null;comment:日志标题" json:"data"`
	CreatedAt time.Time `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null;comment:更新时间" json:"updated_at"` // 更新时间
	DeletedAt time.Time `gorm:"column:deleted_at;comment:删除时间" json:"deleted_at"`          // 删除时间
}

func (game *Game) TableName() string {
	return TableNameGame
}
