package container

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"

	"github.com/LittleBenx86/Benlog/internal/global/consts"
)

var (
	concurrentMap sync.Map
	instance      *cacheContainer
)

type cacheContainer struct {
	Logger *zap.Logger
}

func New(l *zap.Logger) {
	instance = &cacheContainer{
		Logger: l,
	}
}

func GetInstance() *cacheContainer {
	return instance
}

func (c *cacheContainer) IsExists(key string) (interface{}, bool) {
	return concurrentMap.Load(key)
}

func (c *cacheContainer) Set(key string, value interface{}) (result bool) {
	if _, exists := c.IsExists(key); exists == false {
		concurrentMap.Store(key, value)
		result = true
	} else {
		msg := fmt.Sprintf("%s, details: %s", consts.ERRORS_CACHE_CONTAINER_DUPLICATED_KEYS, "key --> "+key)
		// Benlog bootstrap step, the zaplog no load, we have to use the go log to print boot logs
		if c.Logger == nil {
			log.Fatal(msg)
		} else {
			c.Logger.Warn(msg)
		}
	}
	return
}

func (c *cacheContainer) Get(key string) interface{} {
	if val, exists := c.IsExists(key); exists {
		return val
	}
	return nil
}

func (c *cacheContainer) Remove(key string) {
	concurrentMap.Delete(key)
}

func (c *cacheContainer) FuzzyRemoveByPrefix(keyPrefix string) {
	concurrentMap.Range(func(k, v interface{}) bool {
		if item, ok := k.(string); ok {
			if strings.HasPrefix(item, keyPrefix) {
				concurrentMap.Delete(item)
			}
		}
		return true
	})
}
