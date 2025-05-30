package model

type User struct {
	Uid          uint64 `json:"uid" gorm:"autoIncrement;primaryKey"`
	Username     string `json:"username" gorm:"unique"`
	Password     string `json:"password"`
	CorePassword string `json:"core_password"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
