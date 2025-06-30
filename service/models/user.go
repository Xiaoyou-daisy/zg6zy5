package models

type Users struct {
	Id       int64  `gorm:"column:id;type:int;comment:账户ID;primaryKey;" json:"id"`           // 账户ID
	Mobile   string `gorm:"column:mobile;type:char(11);comment:手机号;not null;" json:"mobile"` // 手机号
	Password string `gorm:"column:password;type:char(32);comment:密码;" json:"password"`       // 密码
}
