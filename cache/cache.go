package cache

import (
	"github.com/jiandan100/gocommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	Redis = "redis"
)

var cacheCfg *CheConfig

/*type CheConfig struct {
	Redis *RedisConfig `yaml:"redis"`
}*/
type CheConfig struct {
	Redis map[string]*RedisConfig `yaml:"redis"`
}

func Init(cfgFile string) {
	buf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Warn(cfgFile + "文件读取失败")
	}
	err = yaml.Unmarshal(buf, &cacheCfg)
	if err != nil {
		log.Warn(cfgFile + "解析失败")
	}
	redisCfgGroup = cacheCfg.Redis
	if redisCfgGroup != nil {
		redisConnGroup()
	}
}
