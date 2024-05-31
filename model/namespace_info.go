package model

import (
	"nacos/definition"
	"nacos/util"
)

type NamespaceInfo struct {
	ID            uint                 `gorm:"column:id;type:bigint(20);not null;autoIncrement;primaryKey;comment:id"`
	NamespaceID   string               `gorm:"column:namespace_id;type:varchar(128);not null;default:'';uniqueIndex:uk_namespace_id;comment:namespace_id"`
	NamespaceName string               `gorm:"column:namespace_name;type:varchar(128);not null;comment:namespace_name"`
	Description   string               `gorm:"column:description;type:varchar(256);null;comment:description"`
	GmtCreate     *definition.DateTime `gorm:"column:gmt_create;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:gmt_create"`
	GmtModified   *definition.DateTime `gorm:"column:gmt_modified;type:datetime;not null;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:gmt_modified"`
}

type SearchNamespace struct {
	NamespaceID *string `form:"namespaceId"`
}

func (searchNamespace *SearchNamespace) SetNamespaceID() {
	namespaceID := searchNamespace.NamespaceID
	if util.ZeroValue(namespaceID) {
		searchNamespace.NamespaceID = new(string)
	}
}

type AddNamespace struct {
	NamespaceID   string `form:"customNamespaceId" json:"customNamespaceId"`
	NamespaceName string `form:"namespaceName" json:"namespaceName" binding:"required" msg:"namespaceName不能为空"`
	Description   string `form:"namespaceDesc" json:"namespaceDesc" binding:"required" msg:"namespaceDesc不能为空"`
}

type UpdateNamespace struct {
	NamespaceID   string `form:"namespace" json:"namespace" binding:"required" msg:"namespace不能为空"`
	NamespaceName string `form:"namespaceShowName" json:"namespaceShowName" binding:"required" msg:"namespaceShowName不能为空"`
	Description   string `form:"namespaceDesc" json:"namespaceDesc" binding:"required" msg:"namespaceDesc不能为空"`
}

type NamespaceDetail struct {
	NamespaceID   string `json:"namespace"`
	NamespaceName string `json:"namespaceShowName"`
	Description   string `json:"namespaceDesc"`
	Quota         int    `json:"quota"`
	Type          int    `json:"type"`
	ConfigCount   int    `json:"configCount"`
}
