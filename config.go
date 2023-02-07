package errors

// Config 配置项
type Config struct {
	ErrorConnectionFlag string
}

var defaultCfg = &Config{
	ErrorConnectionFlag: "\nCaused by: ",
}

// SetCfg 设置配置
func SetCfg(c *Config) {
	defaultCfg = c
}
