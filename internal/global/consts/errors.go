package consts

const (
	ERRORS_CONFIG_APP_YML_FILE_NOT_EXISTS    = "could not find the application.yml"
	ERRORS_CONFIG_GORMV2_YML_FILE_NOT_EXISTS = "missing the application-gormv2.yml"
	ERRORS_CONFIG_INIT_ERR                   = "initial the configuration failed"
	ERRORS_CACHE_CONTAINER_DUPLICATED_KEYS   = "duplicated key in cache container"
	ERRORS_DB_DRIVER_UNSUPPORTED             = "unsupported database driver"
	ERRORS_DB_DIALECT_INIT_ERR               = "gorm dialect init failed"
	ERRORS_FILE_OPEN_ERR                     = "open file failed"
	ERRORS_FILE_READ_ERR                     = "read file content failed"
	ERRORS_WS_OPEN_FAILED                    = "websocket on open stage error"
	ERRORS_WS_PROTOCOL_UPGRADE_FAILED        = "websocket protocol upgrade failed"
	ERRORS_WS_READ_MESSAGE_FAILED            = "websocket read pump error"
	ERRORS_WS_HEARTBEAT_SERVER_ERROR         = "websocket heart beat goroutine error"
	ERRORS_WS_HEARTBEAT_FAILURE_EXCEED_MAX   = "websocket fail heart beat count exceed max limitation"
	ERRORS_WS_SET_WRITE_DEADLINE_FAILED      = "websocket fail to set write data deadline"
	ERRORS_WS_WRITE_MESSAGE_FAILED           = "websocket file to write/send message"
	ERRORS_WS_STATE_INVALID                  = "websocket state invalid (suspend, offline)"
	ERRORS_EVENT_FN_KEY_DEUPLICATED          = "event function key name duplicated"
	ERRORS_EVENT_FN_CALL_FAILED              = "event call back function execute failed"
	ERRORS_EVENT_FN_UNREGISTER_TO_CALL       = "event call back function unregister"
)
