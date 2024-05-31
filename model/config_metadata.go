package model

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"nacos/util"
	"nacos/util/collection"
	"strings"
)

type configMetadata struct {
	GroupID     string `yaml:"group"`
	DataID      string `yaml:"dataId"`
	Type        string `yaml:"type"`
	AppName     string `yaml:"appName"`
	Description string `yaml:"desc"`
}

type exportConfigMetadata struct {
	Metadata []configMetadata `yaml:"metadata"`
}

type ConfigMetadataHandler interface {
	Parse(metadata string) []ConfigInfo
	Generate(configInfos []ConfigInfo) string
}

func NewConfigMetadataHandler(v1 bool) ConfigMetadataHandler {
	if v1 {
		return &v1ConfigMetadataHandler{}
	} else {
		return &v2ConfigMetadataHandler{}
	}
}

type v1ConfigMetadataHandler struct{}

func (handler *v1ConfigMetadataHandler) Parse(data string) []ConfigInfo {
	var configInfos []ConfigInfo
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		lines := strings.Split(line, "=")
		if len(lines) == 2 {
			keys := strings.Split(lines[0], ".")
			if len(keys) == 3 {
				configInfo := ConfigInfo{}
				configInfo.GroupID = keys[0]
				configInfo.DataID = keys[1]
				configInfo.AppName = lines[1]
				configInfos = append(configInfos, configInfo)
			}
		}
	}
	return configInfos
}

func (handler *v1ConfigMetadataHandler) Generate(configInfos []ConfigInfo) string {
	joiner := collection.NewJoiner("\n", "", "")
	for _, configInfo := range configInfos {
		joiner.Append(fmt.Sprintf("%s.%s.app=%s", configInfo.GroupID, configInfo.DataID, configInfo.AppName))
	}
	return joiner.String()
}

type v2ConfigMetadataHandler struct{}

func (handler *v2ConfigMetadataHandler) Parse(data string) []ConfigInfo {
	var configInfos []ConfigInfo
	metadataArray := exportConfigMetadata{}
	yaml.Unmarshal([]byte(data), &metadataArray)
	for _, metadata := range metadataArray.Metadata {
		configInfo := ConfigInfo{}
		util.Copy(&metadata, &configInfo)
		configInfos = append(configInfos, configInfo)
	}
	return configInfos
}

func (handler *v2ConfigMetadataHandler) Generate(configInfos []ConfigInfo) string {
	exportMetadata := exportConfigMetadata{}
	for _, configInfo := range configInfos {
		metadata := configMetadata{}
		util.Copy(&configInfo, &metadata)
		exportMetadata.Metadata = append(exportMetadata.Metadata, metadata)
	}
	data, _ := yaml.Marshal(&exportMetadata)
	return string(data)
}
