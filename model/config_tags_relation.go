package model

type ConfigTagsRelation struct {
	ID          uint   `gorm:"column:id;type:bigint(20);not null;autoIncrement;primaryKey;comment:id"`
	NamespaceID string `gorm:"column:namespace_id;type:varchar(128);not null;default:'';uniqueIndex:uk_namespace_id_group_id_data_id_tag_name;comment:namespace_id"`
	GroupID     string `gorm:"column:group_id;type:varchar(128);not null;uniqueIndex:uk_namespace_id_group_id_data_id_tag_name;comment:group_id"`
	DataID      string `gorm:"column:data_id;type:varchar(255);not null;uniqueIndex:uk_namespace_id_group_id_data_id_tag_name;comment:data_id"`
	TagName     string `gorm:"column:tag_name;type:varchar(128);not null;uniqueIndex:uk_namespace_id_group_id_data_id_tag_name;comment:tag_name"`
}
