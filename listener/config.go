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
	listeners: map[string][]chan string{},
}

type configListenerManager struct {
	mutex     sync.RWMutex
	listeners map[string][]chan string
}

func (manager *configListenerManager) Add(key string, ch chan string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	manager.listeners[key] = append(manager.listeners[key], ch)
}

func (manager *configListenerManager) Remove(key string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.listeners, key)
}

func (manager *configListenerManager) Notice(configKey model.ConfigKey) {
	key := fmt.Sprintf("%s+%s+%s", *configKey.NamespaceID, configKey.GroupID, configKey.DataID)
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	if channels, ok := manager.listeners[key]; ok {
		for _, channel := range channels {
			channel <- BuildChangedKey(configKey)
		}
	}
}

func BuildChangedKey(configKey model.ConfigKey) string {
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
