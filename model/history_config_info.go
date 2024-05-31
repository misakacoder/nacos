package model

import (
	"nacos/definition"
	"time"
)

type HistoryConfigInfo struct {
	ID          uint                 `gorm:"column:id;type:bigint(20);not null;autoIncrement;primaryKey;comment:id"`
	NamespaceID string               `gorm:"column:namespace_id;type:varchar(128);not null;default:'';index:idx_namespace_id_group_id_data_id;comment:namespace_id"`
	GroupID     string               `gorm:"column:group_id;type:varchar(128);not null;index:idx_namespace_id_group_id_data_id;comment:group_id"`
	DataID      string               `gorm:"column:data_id;type:varchar(255);not null;index:idx_namespace_id_group_id_data_id;comment:data_id"`
	Content     string               `gorm:"column:content;type:longtext;not null;comment:content"`
	MD5         string               `gorm:"column:md5;type:varchar(32);not null;comment:md5"`
	Type        string               `gorm:"column:type;type:varchar(64);not null;comment:type"`
	SrcUser     string               `gorm:"column:src_user;type:text;null;comment:source user"`
	SrcIP       string               `gorm:"column:src_ip;type:varchar(50);null;comment:source ip"`
	AppName     string               `gorm:"column:app_name;type:varchar(128);null;comment:app name"`
	Description string               `gorm:"column:description;type:varchar(256);null;comment:description"`
	OpType      string               `gorm:"column:op_type;type:char(10);null;comment:source user"`
	GmtCreate   *definition.DateTime `gorm:"column:gmt_create;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:gmt_create"`
	GmtModified *definition.DateTime `gorm:"column:gmt_modified;type:datetime;not null;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:gmt_modified"`
}

type GetHistoryConfig struct {
	ConfigKey
	ID uint `form:"nid" binding:"required" msg:"nid不能为空"`
}

type HistoryConfigDetail struct {
	ID          string    `json:"id"`
	NamespaceID string    `json:"tenant"`
	GroupID     string    `json:"group"`
	DataID      string    `json:"dataId"`
	Content     string    `json:"content"`
	MD5         string    `json:"md5"`
	Type        string    `json:"type"`
	SrcUser     string    `json:"srcUser"`
	SrcIP       string    `json:"srcIp"`
	AppName     string    `json:"appName"`
	Description string    `json:"description"`
	OpType      string    `json:"opType"`
	GmtCreate   time.Time `json:"createdTime"`
	GmtModified time.Time `json:"lastModifiedTime"`
}
