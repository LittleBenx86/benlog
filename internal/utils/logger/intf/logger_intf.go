package intf

import "go.uber.org/zap"

type LoggerContext interface {
	GetLogger() *zap.Logger
}
