package intf

import "time"

type YmlConfig interface {
	ConfigFileUpdateListen()
	Clone(fileName string) YmlConfig
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
}

type YmlContext interface {
	GetYml() YmlConfig
}
