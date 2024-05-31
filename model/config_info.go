package model

import (
	"nacos/definition"
	"nacos/util"
	"nacos/util/collection"
	"strconv"
	"strings"
)

type ConfigInfo struct {
	ID          uint                 `gorm:"column:id;type:bigint(20);not null;autoIncrement;primaryKey;comment:id"`
	NamespaceID string               `gorm:"column:namespace_id;type:varchar(128);not null;default:'';uniqueIndex:uk_namespace_id_group_id_data_id;comment:namespace_id"`
	GroupID     string               `gorm:"column:group_id;type:varchar(128);not null;uniqueIndex:uk_namespace_id_group_id_data_id;comment:group_id"`
	DataID      string               `gorm:"column:data_id;type:varchar(255);not null;uniqueIndex:uk_namespace_id_group_id_data_id;comment:data_id"`
	Content     string               `gorm:"column:content;type:longtext;not null;comment:content"`
	MD5         string               `gorm:"column:md5;type:varchar(32);not null;comment:md5"`
	Type        string               `gorm:"column:type;type:varchar(64);not null;comment:type"`
	AppName     string               `gorm:"column:app_name;type:varchar(128);null;comment:app name"`
	Description string               `gorm:"column:description;type:varchar(256);null;comment:description"`
	SrcUser     string               `gorm:"column:src_user;type:text;null;comment:source user"`
	SrcIP       string               `gorm:"column:src_ip;type:varchar(50);null;comment:source ip"`
	GmtCreate   *definition.DateTime `gorm:"column:gmt_create;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:gmt_create"`
	GmtModified *definition.DateTime `gorm:"column:gmt_modified;type:datetime;not null;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:gmt_modified"`
}

type ConfigKey struct {
	NamespaceID *string `form:"tenant" json:"tenant"`
	GroupID     string  `form:"group" json:"group" binding:"required" msg:"group不能为空"`
	DataID      string  `form:"dataId" json:"dataId" binding:"required" msg:"dataId不能为空"`
}

func (configKey *ConfigKey) SetNamespaceID() {
	if util.ZeroValue(configKey.NamespaceID) {
		configKey.NamespaceID = new(string)
	}
}

type SearchConfig struct {
	NamespaceID string `form:"tenant"`
	GroupID     string `form:"group"`
	DataID      string `form:"dataId"`
	AppName     string `form:"appName"`
	ConfigTags  string `form:"config_tags"`
	Content     string `form:"config_detail"`
}

type AddConfig struct {
	ConfigKey
	Content     string `form:"content" json:"content" binding:"required" msg:"content不能为空"`
	Type        string `form:"type" json:"type"`
	Description string `form:"desc" json:"desc"`
	AppName     string `form:"appName" json:"appName"`
	ConfigTags  string `form:"config_tags" json:"config_tags"`
}

type SelectConfig struct {
	NamespaceID string `form:"tenant" json:"tenant"`
	ID          string `form:"ids" json:"ids"`
}

func (selectConfig *SelectConfig) GetIDArray() []int {
	var result []int
	ids := collection.Distinct(strings.Split(selectConfig.ID, ","))
	for _, id := range ids {
		arg, err := strconv.Atoi(id)
		if err == nil {
			result = append(result, arg)
		}
	}
	return result
}

type ListenerConfig struct {
	LongPullingTimeout string `header:"Long-Pulling-Timeout"`
	ListeningConfigs   string `form:"Listening-Configs" binding:"required" msg:"Listening-Configs不能为空"`
}

type CloneConfigParam struct {
	NamespaceID string `form:"tenant"`
	Policy      string `form:"policy" binding:"required" msg:"policy不能为空"`
}

type CloneConfigBody struct {
	ID      string `json:"cfgId" binding:"required" msg:"cfgId不能为空"`
	GroupID string `json:"group" binding:"required" msg:"group不能为空"`
	DataID  string `json:"dataId" binding:"required" msg:"dataId不能为空"`
}

type ConfigDetail struct {
	ID          string `json:"id"`
	NamespaceID string `json:"tenant"`
	GroupID     string `json:"group"`
	DataID      string `json:"dataId"`
	Content     string `json:"content"`
	MD5         string `json:"md5"`
	Type        string `json:"type"`
	AppName     string `json:"appName"`
	Description string `json:"desc"`
	SrcUser     string `json:"createUser"`
	SrcIP       string `json:"createIp"`
	ConfigTags  string `json:"configTags"`
	CreateTime  int64  `json:"createTime"`
	ModifyTime  int64  `json:"modifyTime"`
}

type BatchAddConfigResultData struct {
	DataId string `json:"dataId"`
	Group  string `json:"group"`
}

type BatchAddConfigResult struct {
	SuccessCount int                        `json:"succCount"`
	SkipCount    int                        `json:"skipCount"`
	SkipData     []BatchAddConfigResultData `json:"skipData"`
	FailData     []BatchAddConfigResultData `json:"failData"`
}
