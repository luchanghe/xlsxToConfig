package Config

import (
	"flag"
)

type Config struct {
	Input  string
	OutPut string
	Files  string
}

var confTemp *Config

func GetConfig() *Config {
	if confTemp != nil {
		return confTemp
	}
	c := &Config{}
	// 使用 flag 包来定义命令行参数
	flag.StringVar(&c.OutPut, "o", "", "./")
	flag.StringVar(&c.Input, "i", "", "./")
	flag.StringVar(&c.Files, "f", "", "")
	// 解析命令行参数
	flag.Parse()
	// 获取字符串的长度
	if c.OutPut[len(c.OutPut)-1] != '/' {
		c.OutPut = c.OutPut + "/"
	}
	if c.OutPut[len(c.OutPut)-1] != '/' {
		c.OutPut = c.OutPut + "/"
	}
	confTemp = c
	return c
}
