package errors

import (
	"sync"
)

// Config represents the configuration options.
type Config struct {
	StackDepth          int    // StackDepth specifies the depth of the function call stack trace. Default value is 10.
	ErrorConnectionFlag string // ErrorConnectionFlag specifies the error connection flag string. Default value is "\nCaused by: ".
}

var (
	cfg        = defaultCfg
	defaultCfg = &Config{
		StackDepth:          10,
		ErrorConnectionFlag: "\nCaused by: ",
	}
	rw sync.RWMutex
)

// SetCfg sets the global configuration instance.
func SetCfg(c *Config) {
	if c == nil {
		return
	}
	rw.Lock()
	defer rw.Unlock()
	cfg = c
}

// GetCfg retrieves the global configuration instance.
func GetCfg() *Config {
	rw.RLock()
	defer rw.RUnlock()
	return cfg
}

// ResetCfg resets the global configuration to its default value.
func ResetCfg() {
	rw.Lock()
	defer rw.Unlock()
	cfg = defaultCfg
}
