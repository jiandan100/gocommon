package db

import (
	"github.com/jiandan100/gocommon/log"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"xorm.io/xorm"
)

var (
	dbCfg     *dbConfig
	DataGroup map[string]*xorm.EngineGroup
)

type dbGroupConfig struct {
	OpenConns       int `yaml:"openConns"`
	IdleConns       int `yaml:"idleConns"`
	ConnMaxLifetime int `yaml:"maxLifetime"`

	Master string   `yaml:"master"`
	Slaves []string `yaml:"slaves"`
}
type dbConfig struct {
	Adapter string                    `yaml:"adapter"`
	Db      map[string]*dbGroupConfig `yaml:"db"`
}

func initDataGroup() map[string]*xorm.EngineGroup {
	var groups = make(map[string]*xorm.EngineGroup)
	if dbCfg == nil {
		log.Error("db config setting error")
	}
	for g, e := range dbCfg.Db {
		dataSourceSlice := make([]string, 0)
		dataSourceSlice = append(dataSourceSlice, e.Master)
		for _, sn := range dbCfg.Db[g].Slaves {
			dataSourceSlice = append(dataSourceSlice, sn)
		}
		if len(dataSourceSlice) > 0 {
			group, err := xorm.NewEngineGroup(dbCfg.Adapter, dataSourceSlice)
			if err != nil {
				log.Warn("创建数据组链接错误：" + err.Error())
			}
			group.SetMaxOpenConns(dbCfg.Db[g].OpenConns)
			group.SetMaxIdleConns(dbCfg.Db[g].IdleConns)
			group.SetConnMaxLifetime(time.Duration(dbCfg.Db[g].ConnMaxLifetime) * time.Second)
			//group.SetConnMaxLifetime(5*time.Minute)
			groups[g] = group
			log.Info(fmt.Sprintf("%s EngineGroup Opened", g))
		}
	}
	return groups
}

func Use(dbName string) *xorm.EngineGroup {
	if DataGroup == nil {
		DataGroup = initDataGroup()
	}
	if g, ok := DataGroup[dbName]; ok {
		return g
	} else {
		log.Error(dbName + " - Database does not exist.")
	}
	return nil
}

func Init(dbCfgFile string) {
	buf, err := ioutil.ReadFile(dbCfgFile)
	if err != nil {
		log.Warn(dbCfgFile + "文件读取失败")
	}
	err = yaml.Unmarshal(buf, &dbCfg)
	if err != nil {
		log.Warn(dbCfgFile + "解析失败")
	}
	DataGroup = initDataGroup()
}

func Close() {
	for n, db := range DataGroup {
		db.Close()
		log.Info(fmt.Sprintf("%s EngineGroup Closed", n))
	}
}
