package cache

import (
	"github.com/jiandan100/gocommon/log"
	"github.com/go-redis/redis"
	"net"
	"time"
)

var (
	RCacheGroup   map[string]*redis.Client
	RCache        *redis.Client
	redisCfgGroup map[string]*RedisConfig
)

type RedisMaster struct {
	Protocol string `yaml:"protocol"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type RedisConfig struct {
	Cluster  string       `yaml:"cluster"`
	Master   *RedisMaster `yaml:"master"`
	Sentinel []string     `yaml:"sentinel"`
}

func redisConnGroup() {
	RCacheGroup = make(map[string]*redis.Client)
	for g, redisCfg := range redisCfgGroup {
		if redisCfg.Master != nil {
			master := redisCfg.Master
			switch redisCfg.Cluster {
			case "", "standalone":
				group := redis.NewClient(&redis.Options{
					Network:  master.Protocol,
					Addr:     net.JoinHostPort(master.Host, master.Port),
					Password: master.Password,
					DB:       master.Db,
				})
				if g == "default" {
					RCache = group
				}
				RCacheGroup[g] = group
			case "sentinel":
				if redisCfg.Sentinel != nil {
					group := redis.NewFailoverClient(&redis.FailoverOptions{
						MasterName:    master.Host,
						Password:      master.Password,
						DB:            master.Db,
						SentinelAddrs: redisCfg.Sentinel,
					})
					if g == "default" {
						RCache = group
					}
					RCacheGroup[g] = group
				} else {
					log.Error("cache config setting error")
				}
			default:
				log.Error("cache config - cluster setting error")
			}
		} else {
			log.Error("cache config - master must be setting")
		}
	}

	if len(RCacheGroup) > 0 {
		log.Info("check connection ... 10s ")
		time.Sleep(10 * time.Second)
		for name, cc := range RCacheGroup {
			_, err := cc.Ping().Result()
			if err != nil {
				log.Error(err)
			} else {
				log.Info("successful connection to redis-server:" + name)
			}
		}
	}
}
