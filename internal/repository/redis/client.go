package redis

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type PoolConfiguration struct {
	Address            string `json:"Address" mapstructure:"address"`
	Username           string `json:"Username" mapstructure:"username"`
	Password           string `json:"Password" mapstructure:"password"`
	DB                 int    `json:"DBIndex" mapstructure:"dbindex"`
	PoolSize           uint   `json:"PoolSize" mapstructure:"poolsize"`
	DialTimeout        uint   `json:"DialTimeout" mapstructure:"dialtimeout"`
	ReadTimeout        int    `json:"ReadTimeout" mapstructure:"readtimeout"`
	WriteTimeout       int    `json:"WriteTimeout" mapstructure:"writetimeout"`
	PoolTimeout        uint   `json:"PoolTimeout" mapstructure:"pooltimeout"`
	IdleTimeout        uint   `json:"IdleTimeout" mapstructure:"idletimeout"`
	MaxRetries         int    `json:"MaxRetries" mapstructure:"maxretries"`
	MinRetryBackoff    int    `json:"MinRetryBackoff" mapstructure:"minretrybackoff"`
	MaxRetryBackoff    int    `json:"MaxRetryBackoff" mapstructure:"maxretrybackoff"`
	IdleCheckFrequency uint   `json:"IdleCheckFrequency" mapstructure:"idlecheckfrequency"`
}

func (r *PoolConfiguration) AsOptions() (opt *redis.Options, err error) {

	if len(r.Address) == 0 {
		err = errors.New("empty address to initLock Redis connection")
		return nil, err
	}

	opt = new(redis.Options)
	opt.Addr = r.Address
	opt.Username = r.Username
	opt.Password = r.Password
	opt.DB = r.DB
	opt.PoolSize = int(r.PoolSize)

	opt.DialTimeout = time.Duration(int64(r.DialTimeout))
	opt.IdleTimeout = time.Duration(int64(r.IdleTimeout))
	opt.PoolTimeout = time.Duration(int64(r.PoolTimeout))
	if r.ReadTimeout < 0 {
		r.ReadTimeout = -1
	}
	opt.ReadTimeout = time.Duration(int64(r.ReadTimeout))

	if r.WriteTimeout < 0 {
		r.WriteTimeout = -1
	}
	opt.WriteTimeout = time.Duration(int64(r.WriteTimeout))

	if r.MaxRetries < 0 {
		r.MaxRetryBackoff = -1
	}
	opt.MaxRetries = r.MaxRetries

	if r.MinRetryBackoff < 0 {
		r.MinRetryBackoff = -1
	}
	opt.MinRetryBackoff = time.Duration(int64(r.MinRetryBackoff))

	if r.MaxRetryBackoff < 0 {
		r.MaxRetryBackoff = -1
	}
	opt.MaxRetryBackoff = time.Duration(int64(r.MaxRetryBackoff))
	opt.IdleCheckFrequency = time.Duration(int64(r.IdleCheckFrequency))
	return
}

type PoolContext struct {
	Cfg *PoolConfiguration
}

var (
	instance *redis.Client
	once     sync.Once
	poolCtx  *PoolContext
)

func NewRedisContext(cfg PoolConfiguration) {
	poolCtx = &PoolContext{
		Cfg: &cfg,
	}
}

// GetInstance
// Use the singleton mode
func GetInstance() *redis.Client {
	once.Do(newRedisClient)
	// After redis client has been initialized, we can pass the nil parameters to get this instance.
	// It means that redis client must be initialized after yaml file was read.
	return instance
}

// newRedisClient
// We only use the single machine redis client
func newRedisClient() {
	opt, err := poolCtx.Cfg.AsOptions()
	if err != nil {
		log.Fatalf("redis configuration convert to option failed! reason: %s\n", err.Error())
	}

	instance = redis.NewClient(opt)
}
