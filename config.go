package errors

// Config 配置项
type Config struct {
	NoStack  bool
	NoDetail bool
}

var defaultCfg = &Config{
	NoStack:  false,
	NoDetail: false,
}

// SetCfg 设置配置
func SetCfg(c *Config) {
	defaultCfg = c
}
