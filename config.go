package errors

import (
	"sync"
)

// Config 配置项
type Config struct {
	Depth               int
	ErrorConnectionFlag string
}

var (
	cfg        = defaultCfg
	defaultCfg = &Config{
		Depth:               10,
		ErrorConnectionFlag: "\nCaused by: ",
	}
	rw sync.RWMutex
)

// SetCfg 设置配置
func SetCfg(c *Config) {
	if c == nil {
		return
	}
	rw.Lock()
	defer rw.Unlock()
	cfg = c
}

// GetCfg 设置配置
func GetCfg() *Config {
	rw.RLock()
	defer rw.RUnlock()
	return cfg
}

// ResetCfg 重置配置
func ResetCfg() {
	rw.Lock()
	defer rw.Unlock()
	cfg = defaultCfg
}
