package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nacos/listener"
	"nacos/model"
	"nacos/router"
	"nacos/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func searchListenerByIP(context *gin.Context) {
	status := map[string]string{}
	result := map[string]any{
		"collectStatus":           200,
		"lisentersGroupkeyStatus": status,
	}
	ip := context.Query("ip")
	if ip != "" {
		namespaceID := context.Query("namespaceId")
		for key, listeners := range listener.ConfigListenerManager.Listeners {
			keys := strings.Split(key, "+")
			if keys[0] == namespaceID {
				groupID := keys[1]
				dataID := keys[2]
				for _, value := range listeners {
					if value.IP == ip {
						status[fmt.Sprintf("%s+%s", groupID, dataID)] = value.MD5
					}
				}
			}
		}
	}
	context.JSON(http.StatusOK, result)
}

func searchListenerByKey(context *gin.Context) {
	status := map[string]string{}
	result := map[string]any{
		"collectStatus":           200,
		"lisentersGroupkeyStatus": status,
	}
	configKey := model.Bind(context, &model.ConfigKey{})
	configKey.SetNamespaceID()
	key := fmt.Sprintf("%s+%s+%s", *configKey.NamespaceID, configKey.GroupID, configKey.DataID)
	if listeners, ok := listener.ConfigListenerManager.Listeners[key]; ok {
		for _, value := range listeners {
			status[value.IP] = value.MD5
		}
	}
	context.JSON(http.StatusOK, result)
}

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
	for key, md5 := range configKeys {
		listener.ConfigListenerManager.Add(key, listener.ConfigListener{IP: util.GetClientIP(context), MD5: md5, Channel: ch})
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
