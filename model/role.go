package model

type Role struct {
	ID       uint   `gorm:"column:id;type:bigint(20);not null;primaryKey;autoIncrement;comment:id"`
	Username string `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username_role,sort:asc;comment:username"`
	Role     string `gorm:"column:role;type:varchar(50);not null;uniqueIndex:uk_username_role,sort:asc;comment:role"`
}

type RoleInfo struct {
	Username string `form:"username" json:"username" binding:"required" msg:"username不能为空"`
	Role     string `form:"role" json:"role" binding:"required" msg:"role不能为空"`
}

type SearchRole struct {
	Username   string `form:"username"`
	Role       string `form:"role"`
	SearchType string `form:"search"`
}

type RoleDetail struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
