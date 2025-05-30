package model

type Account struct {
	Rid        uint64 `json:"rid" gorm:"primaryKey"`
	Uid        uint64 `json:"uid" gorm:"not null;index:idx_uid_update_time,priority:1"`
	Website    string `json:"website" gorm:"type:varchar(2048);not null"`
	Account    string `json:"account" gorm:"type:varchar(2048)"`
	Password   string `json:"password" gorm:"type:varchar(4096);not null"`
	UpdateTime int64  `json:"update_time" gorm:"not null;index:idx_uid_update_time,priority:2"`
}

// TableName 指定表名
func (Account) TableName() string {
	return "accounts"
}
