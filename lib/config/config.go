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
	flag.StringVar(&c.OutPut, "o", "", "./")
	flag.StringVar(&c.Input, "i", "", "./")
	flag.StringVar(&c.Files, "f", "", "")
	flag.Parse()
	if c.OutPut[len(c.OutPut)-1] != '/' {
		c.OutPut = c.OutPut + "/"
	}
	if c.OutPut[len(c.OutPut)-1] != '/' {
		c.OutPut = c.OutPut + "/"
	}
	confTemp = c
	return c
}
