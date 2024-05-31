package model

import (
	"fmt"
	"nacos/consts"
	"nacos/database"
	"nacos/util"
)

const tableOptionsFormat = "engine=InnoDB default charset=utf8mb4 collate=utf8mb4_bin comment='%s'"

func FirstOrCreate() {
	tables := map[string]any{
		"config_info":          &ConfigInfo{},
		"history_config_info":  &HistoryConfigInfo{},
		"config_tags_relation": &ConfigTagsRelation{},
		"namespace_info":       &NamespaceInfo{},
		"user":                 &User{},
		"role":                 &Role{},
	}
	for k, v := range tables {
		_ = db.GORM.Set("gorm:table_options", fmt.Sprintf(tableOptionsFormat, k)).AutoMigrate(v)
	}
	init := []any{
		&NamespaceInfo{NamespaceID: "", NamespaceName: consts.DefaultNamespaceID, Description: "Public Namespace"},
		&User{Username: "nacos", Password: util.MD5("nacos"), Enabled: true},
		&Role{Username: "nacos", Role: "ROLE_ADMIN"},
	}
	for _, data := range init {
		db.GORM.FirstOrCreate(data)
	}
}
