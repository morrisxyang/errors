package errors

// Config 配置项
type Config struct {
	NoStack             bool
	NoDetail            bool
	ErrorConnectionFlag string
}

var defaultCfg = &Config{
	NoStack:             false,
	NoDetail:            false,
	ErrorConnectionFlag: "\nCaused by: ",
}

// SetCfg 设置配置
func SetCfg(c *Config) {
	defaultCfg = c
}
