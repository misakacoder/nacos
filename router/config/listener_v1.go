package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nacos/model"
	"nacos/router"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mutex           = sync.Mutex{}
	configListeners = map[string][]chan string{}
)

func listenConfig(context *gin.Context) {
	listenerConfig := model.Bind(context, &model.ListenerConfig{})
	model.BindHeader(context, listenerConfig)
	listeningConfigs, _ := url.QueryUnescape(listenerConfig.ListeningConfigs)
	pullingTimeout, err := strconv.Atoi(listenerConfig.LongPullingTimeout)
	if err != nil {
		pullingTimeout = defaultPullingTimeout
	}
	var configKey string
	var configKeySlice []string
	configKeys := map[string]string{}
	for _, char := range listeningConfigs {
		if char == 1 {
			configKeySlice = append(configKeySlice, configKey)
			configSliceLength := len(configKeySlice)
			if configSliceLength >= 3 && configSliceLength <= 4 {
				dataID := configKeySlice[0]
				groupID := configKeySlice[1]
				md5 := configKeySlice[2]
				namespaceID := ""
				if configSliceLength == 4 {
					namespaceID = configKeySlice[3]
				}
				key := fmt.Sprintf("%s+%s+%s", namespaceID, groupID, dataID)
				configKeys[key] = md5
			} else {
				router.ConfigListenerError.With("监听的参数格式错误").Error(context)
				return
			}
			configKeySlice = []string{}
			configKey = ""
		} else if char == 2 {
			configKeySlice = append(configKeySlice, configKey)
			configKey = ""
		} else {
			configKey += string(char)
		}
	}
	for key, md5 := range configKeys {
		keys := strings.Split(key, "+")
		baseConfigRO := model.ConfigKey{NamespaceID: &keys[0], GroupID: keys[1], DataID: keys[2]}
		configInfo := getConfigInfo(context, &baseConfigRO)
		if configInfo != nil && configInfo.MD5 != md5 {
			context.String(http.StatusOK, buildChangedConfigInfo(baseConfigRO))
			return
		}
	}
	ch := make(chan string)
	defer close(ch)
	defer func() {
		mutex.Lock()
		defer mutex.Unlock()
		for key := range configKeys {
			delete(configListeners, key)
		}
	}()
	func() {
		mutex.Lock()
		defer mutex.Unlock()
		for key := range configKeys {
			configListeners[key] = append(configListeners[key], ch)
		}
	}()
	timer := time.NewTimer(time.Duration(pullingTimeout) * time.Millisecond)
	select {
	case changeKey := <-ch:
		context.String(http.StatusOK, changeKey)
		break
	case <-timer.C:
		break
	}
}

func buildChangedConfigInfo(configKey model.ConfigKey) string {
	namespaceID := *configKey.NamespaceID
	builder := strings.Builder{}
	builder.WriteString(configKey.DataID)
	builder.WriteRune(2)
	builder.WriteString(configKey.GroupID)
	if namespaceID == "" {
		builder.WriteRune(1)
	} else {
		builder.WriteRune(2)
		builder.WriteString(namespaceID)
		builder.WriteRune(1)
	}
	return url.QueryEscape(builder.String())
}
