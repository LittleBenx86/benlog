package core

import (
	"fmt"
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	"strings"
	"sync"
)

var (
	eventMgrInstance *EventManager
	initLock         sync.Mutex
	SysEventMap      sync.Map
)

type EventManager struct {
}

func GetEventMgrInstance() *EventManager {
	if eventMgrInstance == nil {
		initLock.Lock()
		defer initLock.Unlock()

		eventMgrInstance = &EventManager{}
	}
	return eventMgrInstance
}

func (e *EventManager) Get(key string) (interface{}, bool) {
	if value, exists := SysEventMap.Load(key); exists {
		return value, exists
	}
	return nil, false
}

func (e *EventManager) Set(key string, keyFn func(args ...interface{})) bool {
	if _, exists := e.Get(key); !exists {
		SysEventMap.Store(key, keyFn)
		return true
	}
	msg := fmt.Sprintf("%s, key name [%s]", consts.ERRORS_EVENT_FN_KEY_DEUPLICATED, key)
	variables.Logger.Warn(msg)
	return false
}

func (e *EventManager) Remove(key string) {
	SysEventMap.Delete(key)
}

func (e *EventManager) Exec(key string, args ...interface{}) {
	valueIntf, exists := e.Get(key)
	if !exists {
		msg := fmt.Sprintf("%s, key fn name [%s]", consts.ERRORS_EVENT_FN_UNREGISTER_TO_CALL, key)
		variables.Logger.Error(msg)
		return
	}

	fn, ok := valueIntf.(func(args ...interface{}))
	if !ok {
		msg := fmt.Sprintf("%s, key fn name [%s]", consts.ERRORS_EVENT_FN_CALL_FAILED, key)
		variables.Logger.Error(msg)
		return
	}

	fn(args...) // real execute event callback function
}

func (e *EventManager) FuzzyExec(keyPrefix string) {
	SysEventMap.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok {
			if strings.HasPrefix(key, keyPrefix) {
				e.Exec(key)
			}
		}
		return false // not real meaning
	})
}
