package errors

// Config 配置项
type Config struct {
	Depth               int
	ErrorConnectionFlag string
}

var defaultCfg = &Config{
	Depth:               10,
	ErrorConnectionFlag: "\nCaused by: ",
}

// SetCfg 设置配置
func SetCfg(c *Config) {
	defaultCfg = c
}
