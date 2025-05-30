package model

type Timestamp struct {
	Uid              uint64 `json:"uid" gorm:"primaryKey"`
	LatestUpdateTime int64  `json:"latest_update_time" gorm:"not null;default:0"`
	LatestDeleteTime int64  `json:"latest_delete_time" gorm:"not null;default:0"`
}

// TableName 指定表名
func (Timestamp) TableName() string {
	return "timestamps"
}
