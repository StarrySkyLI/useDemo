package schema

import "time"

const TableNameFile = "file"

// File mapped from table <file>
type File struct {
	ID         int64     `gorm:"column:id;type:bigint;primaryKey" json:"id"`
	DomainName string    `gorm:"column:domain_name;type:varchar(100);not null" json:"domain_name"`
	BizName    string    `gorm:"column:biz_name;type:varchar(100);not null" json:"biz_name"`
	Hash       string    `gorm:"column:hash;type:varchar(255);not null;index:hash_idx,priority:1" json:"hash"`
	FileSize   int64     `gorm:"column:file_size;type:bigint;not null" json:"file_size"`
	FileType   string    `gorm:"column:file_type;type:varchar(255);not null" json:"file_type"`
	Uploaded   bool      `gorm:"column:uploaded;type:tinyint(1);not null" json:"uploaded"`
	IsDeleted  bool      `gorm:"column:is_deleted;type:tinyint(1);not null" json:"is_deleted"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null;index:create_time_idx,priority:1;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;not null;index:update_time_idx,priority:1;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName File's table name
func (*File) TableName() string {
	return TableNameFile
}
