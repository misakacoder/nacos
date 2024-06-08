package listener

import (
	"fmt"
	"nacos/model"
	"net/url"
	"strings"
	"sync"
)

var ConfigListenerManager = &configListenerManager{
	mutex:     sync.RWMutex{},
	Listeners: map[string][]ConfigListener{},
}

type ConfigListener struct {
	IP      string
	MD5     string
	Channel chan string
}

type configListenerManager struct {
	mutex     sync.RWMutex
	Listeners map[string][]ConfigListener
}

func (manager *configListenerManager) Add(key string, listener ConfigListener) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.Listeners[key] = append(manager.Listeners[key], listener)
}

func (manager *configListenerManager) Remove(key string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.Listeners, key)
}

func (manager *configListenerManager) Notify(configKey model.ConfigKey) {
	namespaceID := ""
	if configKey.NamespaceID != nil {
		namespaceID = *configKey.NamespaceID
	}
	key := fmt.Sprintf("%s+%s+%s", namespaceID, configKey.GroupID, configKey.DataID)
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	if listeners, ok := manager.Listeners[key]; ok {
		for _, listener := range listeners {
			go func(channel chan string) {
				channel <- BuildChangedKey(configKey)
			}(listener.Channel)
		}
	}
}

func BuildChangedKey(configKey model.ConfigKey) string {
	namespaceID := ""
	if configKey.NamespaceID != nil {
		namespaceID = *configKey.NamespaceID
	}
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
