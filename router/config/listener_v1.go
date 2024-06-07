package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nacos/listener"
	"nacos/model"
	"nacos/router"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func listenConfig(context *gin.Context) {
	listenerConfig := model.Bind(context, &model.ListenerConfig{})
	model.BindHeader(context, listenerConfig)
	listeningConfigs, _ := url.QueryUnescape(listenerConfig.ListeningConfigs)
	pullingTimeout, err := strconv.Atoi(listenerConfig.LongPullingTimeout)
	if err != nil {
		pullingTimeout = defaultPullingTimeout
	}
	var listeningKey string
	var listeningKeys []string
	configKeys := map[string]string{}
	for _, char := range listeningConfigs {
		if char == 1 {
			listeningKeys = append(listeningKeys, listeningKey)
			configSliceLength := len(listeningKeys)
			if configSliceLength >= 3 && configSliceLength <= 4 {
				dataID := listeningKeys[0]
				groupID := listeningKeys[1]
				md5 := listeningKeys[2]
				namespaceID := ""
				if configSliceLength == 4 {
					namespaceID = listeningKeys[3]
				}
				key := fmt.Sprintf("%s+%s+%s", namespaceID, groupID, dataID)
				configKeys[key] = md5
			} else {
				router.ConfigListenerError.With("监听的参数格式错误").Error(context)
				return
			}
			listeningKeys = []string{}
			listeningKey = ""
		} else if char == 2 {
			listeningKeys = append(listeningKeys, listeningKey)
			listeningKey = ""
		} else {
			listeningKey += string(char)
		}
	}
	for key, md5 := range configKeys {
		keys := strings.Split(key, "+")
		configKey := model.ConfigKey{NamespaceID: &keys[0], GroupID: keys[1], DataID: keys[2]}
		configInfo := getConfigInfo(context, &configKey)
		if configInfo != nil && configInfo.MD5 != md5 {
			context.String(http.StatusOK, listener.BuildChangedKey(configKey))
			return
		}
	}
	ch := make(chan string)
	defer close(ch)
	defer func() {
		for key := range configKeys {
			listener.ConfigListenerManager.Remove(key)
		}
	}()
	for key := range configKeys {
		listener.ConfigListenerManager.Add(key, ch)
	}
	timer := time.NewTimer(time.Duration(pullingTimeout) * time.Millisecond)
	select {
	case changeKey := <-ch:
		context.String(http.StatusOK, changeKey)
		break
	case <-timer.C:
		break
	}
}
