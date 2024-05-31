package model

type Role struct {
	ID       uint   `gorm:"column:id;type:bigint(20);not null;primaryKey;autoIncrement;comment:id"`
	Username string `gorm:"column:username;type:varchar(50);not null;uniqueIndex:uk_username_role,sort:asc;comment:username"`
	Role     string `gorm:"column:role;type:varchar(50);not null;uniqueIndex:uk_username_role,sort:asc;comment:role"`
}
