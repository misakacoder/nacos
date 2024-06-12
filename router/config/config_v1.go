package config

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"nacos/cluster"
	"nacos/database"
	"nacos/database/dbutil"
	"nacos/listener"
	"nacos/model"
	"nacos/router"
	"nacos/util"
	"nacos/util/collection"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPullingTimeout = 60000
	v1Metadata            = ".meta.yml"
	v2Metadata            = ".metadata.yml"
)

var metadataHandler = map[string]model.ConfigMetadataHandler{
	v1Metadata: model.NewConfigMetadataHandler(true),
	v2Metadata: model.NewConfigMetadataHandler(false),
}

func queryConfig(context *gin.Context) {
	if searchType := context.Query("search"); searchType != "" {
		searchConfig(context, searchType)
	} else if showType := context.Query("show"); showType != "" {
		show(context, showType)
	} else if exportV1, exportV2 := context.Query("export"), context.Query("exportV2"); exportV1 != "" || exportV2 != "" {
		v1 := util.ConditionalExpression(exportV1 == "true", true, false)
		export(context, v1)
	} else {
		getConfigContent(context)
	}
}

func searchConfig(context *gin.Context, searchType string) {
	blur := searchType == "blur"
	param := model.Bind(context, &model.SearchConfig{})
	page := model.Bind(context, &model.Page{})
	sql, args := getSqlBuilder(param, blur).Build()
	context.JSON(http.StatusOK, model.PaginateSQL[model.ConfigDetail](sql, args, page))
}

func getSQLBuilder(searchConfig *model.SearchConfig, blur bool) dbutil.QueryBuilder {
	appName := searchConfig.AppName
	sql := `select
    		distinct
			config.id, 
			config.namespace_id, 
			config.group_id, 
			config.data_id, 
			config.content, 
			config.md5, 
			config.type, 
			config.app_name 
			from config_info config 
			left join config_tags_relation tag 
			on config.namespace_id = tag.namespace_id 
			and config.group_id = tag.group_id 
			and config.data_id = tag.data_id`
	sqlBuilder := dbutil.NewSQLBuilder(sql)
	sqlBuilder.Where("config.namespace_id = ?", searchConfig.NamespaceID)
	sqlBuilder.WhereOnConditional("config.app_name = ?", appName, appName != "")
	whereOnBlur(sqlBuilder, "config.group_id", searchConfig.GroupID, blur)
	whereOnBlur(sqlBuilder, "config.data_id", searchConfig.DataID, blur)
	whereOnBlur(sqlBuilder, "config.content", searchConfig.Content, blur)
	configTags := searchConfig.ConfigTags
	if configTags != "" {
		tags := collection.Distinct(strings.Split(configTags, ","))
		sqlBuilder.WhereOnConditional("tag.tag_name in ?", tags, true)
	}
	return sqlBuilder
}

func getSqlBuilder(searchConfig *model.SearchConfig, blur bool) dbutil.QueryBuilder {
	appName := searchConfig.AppName
	sqlBuilder := dbutil.NewSqlBuilder().
		Select("distinct config.id").
		Select("config.namespace_id").
		Select("config.group_id").
		Select("config.data_id").
		Select("config.content").
		Select("config.md5").
		Select("config.type").
		Select("config.app_name").
		From("config_info config").
		LeftJoin("config_tags_relation tag", "config.namespace_id = tag.namespace_id and config.group_id = tag.group_id and config.data_id = tag.data_id").
		Where("config.namespace_id = ?", searchConfig.NamespaceID, true).
		Where("config.app_name = ?", appName, appName != "")
	whereOnBlur(sqlBuilder, "config.group_id", searchConfig.GroupID, blur)
	whereOnBlur(sqlBuilder, "config.data_id", searchConfig.DataID, blur)
	whereOnBlur(sqlBuilder, "config.content", searchConfig.Content, blur)
	configTags := searchConfig.ConfigTags
	if configTags != "" {
		tags := collection.Distinct(strings.Split(configTags, ","))
		sqlBuilder.WhereOnConditional("tag.tag_name in ?", tags, true)
	}
	return sqlBuilder
}

func whereOnBlur(builder dbutil.QueryBuilder, field, value string, blur bool) {
	sql, arg := router.BlurQuery(field, value, blur)
	builder.WhereOnConditional(sql, arg, value != "")
}

func show(context *gin.Context, showType string) {
	if showType == "all" {
		configKey := model.ConfigKey{}
		configInfo := getConfigInfo(context, &configKey)
		if configInfo != nil {
			configDetail := model.ConfigDetail{}
			util.Copy(configInfo, &configDetail)
			configDetail.ID = strconv.Itoa(int(configInfo.ID))
			configDetail.ConfigTags = getConfigTags(configKey)
			configDetail.CreateTime = time.Time(*configInfo.GmtCreate).UnixMilli()
			configDetail.ModifyTime = time.Time(*configInfo.GmtModified).UnixMilli()
			context.JSON(http.StatusOK, configDetail)
		}
	}
}

func export(context *gin.Context, v1 bool) {
	selectConfig := model.Bind(context, &model.SelectConfig{})
	ids := selectConfig.GetIDArray()
	var configInfos []model.ConfigInfo
	conditions := []any{model.ConfigKey{NamespaceID: &selectConfig.NamespaceID}}
	if len(ids) > 0 {
		conditions = append(conditions, []any{"id in ?", ids})
	}
	dbutil.MultiCondition(db.GORM.Model(model.ConfigInfo{}), conditions).Find(&configInfos)
	root, _ := os.MkdirTemp("/", uuid.New().String())
	defer os.RemoveAll(root)
	zipDir, _ := os.MkdirTemp(root, uuid.New().String())
	zipFilename := fmt.Sprintf("nacos_config_export_%s.zip", time.Now().Format("20060102150405"))
	zipFilepath := filepath.Join(root, zipFilename)
	if len(configInfos) > 0 {
		for _, configInfo := range configInfos {
			groupDir := filepath.Join(zipDir, configInfo.GroupID)
			os.MkdirAll(groupDir, 0777)
			dataFile, _ := os.Create(filepath.Join(groupDir, configInfo.DataID))
			dataFile.WriteString(configInfo.Content)
		}
		metadataFilename := util.ConditionalExpression(v1, v1Metadata, v2Metadata)
		metadata := metadataHandler[metadataFilename].Generate(configInfos)
		os.WriteFile(filepath.Join(zipDir, metadataFilename), []byte(metadata), 0777)
		err := util.ZipDir(zipDir, zipFilepath)
		if err != nil {
			router.ServerError.Msg("导出失败").Error(context)
		} else {
			context.FileAttachment(zipFilepath, zipFilename)
		}
	}
}

func getConfigContent(context *gin.Context) {
	content := ""
	configInfo := getConfigInfo(context, &model.ConfigKey{})
	if configInfo != nil {
		content = configInfo.Content
	}
	context.String(http.StatusOK, content)
}

func saveConfig(context *gin.Context) {
	if context.Query("clone") == "true" {
		cloneConfig(context)
	} else if context.Query("import") == "true" {
		importConfig(context)
	} else {
		addConfig(context)
	}
}

func cloneConfig(context *gin.Context) {
	configQuery := model.BindQuery(context, &model.CloneConfigParam{})
	configBodies := model.BindJSON(context, &[]model.CloneConfigBody{})
	namespaceID := configQuery.NamespaceID
	policy := configQuery.Policy
	srcUser := context.GetString("username")
	srcIP := util.GetClientIP(context)
	var configInfos []model.ConfigInfo
	for _, configBody := range *configBodies {
		groupID := configBody.GroupID
		dataID := configBody.DataID
		configInfo := model.ConfigInfo{}
		db.GORM.Where(model.ConfigInfo{ID: util.Atoi[uint](configBody.ID)}).Find(&configInfo)
		configInfo.ID = 0
		configInfo.NamespaceID = namespaceID
		configInfo.GroupID = groupID
		configInfo.DataID = dataID
		configInfo.SrcUser = srcUser
		configInfo.SrcIP = srcIP
		configInfo.GmtCreate = nil
		configInfo.GmtModified = nil
		configInfos = append(configInfos, configInfo)

	}
	result := getBatchAddConfigResult(configInfos, policy)
	router.OK.Msg("Clone Completed Successfully").With(result).Ok(context)
}

func getBatchAddConfigResult(configInfos []model.ConfigInfo, policy string) model.BatchAddConfigResult {
	result := model.BatchAddConfigResult{
		SkipData: []model.BatchAddConfigResultData{},
		FailData: []model.BatchAddConfigResultData{},
	}
loop:
	for i, configInfo := range configInfos {
		namespaceID := configInfo.NamespaceID
		groupID := configInfo.GroupID
		dataID := configInfo.DataID
		if err := db.GORM.Create(&configInfo).Error; errors.Is(err, gorm.ErrDuplicatedKey) {
			switch policy {
			case "ABORT":
				result.FailData = append(result.FailData, model.BatchAddConfigResultData{Group: groupID, DataId: dataID})
				setBatchAddConfigResultSkipData(&result, configInfos[i+1:])
				break loop
			case "SKIP":
				setBatchAddConfigResultSkipData(&result, []model.ConfigInfo{configInfo})
			case "OVERWRITE":
				result.SuccessCount += 1
				db.GORM.Where(&model.ConfigKey{NamespaceID: &namespaceID, GroupID: groupID, DataID: groupID}).Updates(configInfo)
				addHistoryConfigInfo(db.GORM, &configInfo, "U")
			}
		} else {
			result.SuccessCount += 1
			addHistoryConfigInfo(db.GORM, &configInfo, "I")
		}
	}
	return result
}

func setBatchAddConfigResultSkipData(addConfigResult *model.BatchAddConfigResult, configInfo []model.ConfigInfo) {
	var skipData []model.BatchAddConfigResultData
	for _, configBody := range configInfo {
		skipData = append(skipData, model.BatchAddConfigResultData{Group: configBody.GroupID, DataId: configBody.DataID})
	}
	addConfigResult.SkipCount += 1
	addConfigResult.SkipData = append(addConfigResult.SkipData, skipData...)
}

func importConfig(context *gin.Context) {
	namespaceID := context.Query("namespace")
	policy := context.DefaultPostForm("policy", "ABORT")
	srcUser := context.GetString("username")
	srcIP := util.GetClientIP(context)
	fileHeader, _ := context.FormFile("file")
	dir, _ := os.MkdirTemp("/", uuid.New().String())
	defer os.RemoveAll(dir)
	filename := filepath.Join(dir, fileHeader.Filename)
	context.SaveUploadedFile(fileHeader, filename)
	reader, err := zip.OpenReader(filename)
	if err != nil {
		router.ServerError.Msg("导入的文件数据为空").Error(context)
		return
	}
	fileMap := map[string]*zip.File{}
	for _, file := range reader.File {
		fileMap[file.Name] = file
	}
	metadata := fileMap[v1Metadata]
	if metadata == nil {
		metadata = fileMap[v2Metadata]
	}
	if metadata == nil {
		context.String(http.StatusBadRequest, "metadata not found")
		return
	}
	handler := metadataHandler[metadata.FileInfo().Name()]
	rd, _ := metadata.Open()
	defer rd.Close()
	bytes, _ := io.ReadAll(rd)
	if handler != nil {
		configInfos := handler.Parse(string(bytes))
		for i, configInfo := range configInfos {
			data := util.ReadZipFile(fileMap[fmt.Sprintf("%s/%s", configInfo.GroupID, configInfo.DataID)])
			content := string(data)
			configInfos[i].NamespaceID = namespaceID
			configInfos[i].Content = content
			configInfos[i].MD5 = util.MD5(content)
			configInfos[i].SrcUser = srcUser
			configInfos[i].SrcIP = srcIP
		}
		result := getBatchAddConfigResult(configInfos, policy)
		router.OK.Msg("导入成功").With(result).Ok(context)
	}
}

func addConfig(context *gin.Context) {
	param := model.Bind(context, &model.AddConfig{})
	changed := true
	db.Transaction(func(tx *gorm.DB) {
		data := &model.ConfigInfo{}
		util.Copy(param, &data)
		data.MD5 = util.MD5(data.Content)
		data.SrcUser = context.GetString("username")
		data.SrcIP = util.GetClientIP(context)
		configInfo := getConfigInfo(context, &param.ConfigKey)
		if configInfo == nil {
			dbutil.PanicError(tx.Create(data))
			addConfigTagsRelation(tx, param.ConfigKey, param.ConfigTags)
			addHistoryConfigInfo(tx, data, "I")
		} else {
			dbutil.PanicError(tx.Model(&model.ConfigInfo{ID: configInfo.ID}).Updates(data))
			addHistoryConfigInfo(tx, configInfo, "U")
			deleteConfigTagsRelation(tx, param.ConfigKey)
			addConfigTagsRelation(tx, param.ConfigKey, param.ConfigTags)
			if configInfo.MD5 == data.MD5 {
				changed = false
			}
		}
	})
	if changed {
		listener.ConfigListenerManager.Notify(param.ConfigKey)
		cluster.CLUSTER.NotifySlaveConfigListener(param.ConfigKey)
	}
	context.String(http.StatusOK, "%v", true)
}

func delConfig(context *gin.Context) {
	switch context.Query("delType") {
	case "ids":
		deleteConfigByPrimaryKey(context)
	default:
		deleteConfigByUniqueKey(context)
	}
}

func deleteConfigByPrimaryKey(context *gin.Context) {
	selectConfig := model.Bind(context, &model.SelectConfig{})
	ids := selectConfig.GetIDArray()
	var affected int64
	if len(ids) > 0 {
		var keys []model.ConfigKey
		db.GORM.Model(model.ConfigInfo{}).Where("id in ?", ids).Find(&keys)
		affected = deleteConfigByUniqueKeys(context, keys)
	}
	context.JSON(http.StatusOK, router.OK.With(affected > 0))
}

func deleteConfigByUniqueKey(context *gin.Context) {
	configKey := model.Bind(context, &model.ConfigKey{})
	deleteConfigByUniqueKeys(context, []model.ConfigKey{*configKey})
	context.String(http.StatusOK, "%v", true)
}

func deleteConfigByUniqueKeys(context *gin.Context, keys []model.ConfigKey) (allAffected int64) {
	db.Transaction(func(tx *gorm.DB) {
		for _, key := range keys {
			key.SetNamespaceID()
			configInfo := &model.ConfigInfo{}
			if err := tx.First(configInfo, key).Error; err == nil {
				affected := dbutil.PanicError(tx.Delete(model.ConfigInfo{}, key)).RowsAffected
				if affected > 0 {
					configInfo.SrcIP = util.GetClientIP(context)
					addHistoryConfigInfo(tx, configInfo, "D")
					deleteConfigTagsRelation(tx, key)
					listener.ConfigListenerManager.Notify(key)
					cluster.CLUSTER.NotifySlaveConfigListener(key)
					allAffected += affected
				}
			}
		}
	})
	return allAffected
}

func getConfigInfo(context *gin.Context, configKey *model.ConfigKey) *model.ConfigInfo {
	model.Bind(context, configKey)
	configKey.SetNamespaceID()
	data := &model.ConfigInfo{}
	if err := db.GORM.Where(configKey).First(data).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return data
}

func addHistoryConfigInfo(db *gorm.DB, configInfo *model.ConfigInfo, operation string) {
	historyConfigInfo := &model.HistoryConfigInfo{}
	util.Copy(configInfo, historyConfigInfo)
	historyConfigInfo.ID = 0
	historyConfigInfo.OpType = operation
	historyConfigInfo.GmtCreate = nil
	historyConfigInfo.GmtModified = nil
	dbutil.PanicError(db.Create(historyConfigInfo))
}

func getConfigTags(configKey model.ConfigKey) string {
	configKey.SetNamespaceID()
	var tagNames []string
	db.GORM.Model(model.ConfigTagsRelation{}).Select("tag_name").Where(configKey).Find(&tagNames)
	return strings.Join(tagNames, ",")
}

func addConfigTagsRelation(db *gorm.DB, configKey model.ConfigKey, configTags string) {
	if configTags != "" {
		tags := collection.Distinct(strings.Split(configTags, ","))
		configKey.SetNamespaceID()
		configTagsRelation := &model.ConfigTagsRelation{NamespaceID: *configKey.NamespaceID, GroupID: configKey.GroupID, DataID: configKey.DataID}
		for _, tag := range tags {
			configTagsRelation.ID = 0
			configTagsRelation.TagName = tag
			dbutil.PanicError(db.Create(configTagsRelation))
		}
	}
}

func deleteConfigTagsRelation(db *gorm.DB, configKey model.ConfigKey) {
	configKey.SetNamespaceID()
	dbutil.PanicError(db.Delete(model.ConfigTagsRelation{}, configKey))
}
