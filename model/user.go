package model

type User struct {
	ID       uint   `gorm:"column:id;type:bigint(20);not null;autoIncrement;primaryKey;comment:id"`
	Username string `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username;comment:username"`
	Password string `gorm:"column:password;type:varchar(500);not null;comment:password"`
	Enabled  bool   `gorm:"column:enabled;type:boolean;not null;comment:enabled"`
}

type Username struct {
	Username string `form:"username" json:"username" binding:"required" msg:"username不能为空"`
}

type UserInfo struct {
	Username
	Password string `form:"password" json:"password" binding:"required" msg:"password不能为空"`
}

type Token struct {
	AccessToken string `json:"accessToken"`
	Ttl         int64  `json:"tokenTtl"`
	GlobalAdmin bool   `json:"globalAdmin"`
	Username    string `json:"username"`
}

type AccessToken struct {
	AccessToken string `form:"accessToken" header:"accessToken"`
}

type SearchUser struct {
	Username   string `form:"username"`
	SearchType string `form:"search"`
}

type UpdateUser struct {
	Username
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required" msg:"newPassword不能为空"`
}

type UserDetail struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
