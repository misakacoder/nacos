package configuration

import (
	"github.com/jinzhu/configor"
	"github.com/misakacoder/logger"
)

type configuration struct {
	Server struct {
		Bind string
		Port int
	}
	Database struct {
		Host            string
		Port            int
		Username        string
		Password        string
		Name            string
		MaxIdleConn     int    `yaml:"maxIdleConn"`
		MaxOpenConn     int    `yaml:"maxOpenConn"`
		ConnMaxLifeTime string `yaml:"connMaxLifeTime"`
		SlowSqlTime     string `yaml:"slowSqlTime"`
		PrintSql        bool   `yaml:"printSql"`
	}
	Nacos struct {
		Auth struct {
			Enabled    bool
			SecretKey  string `yaml:"secretKey"`
			ExpireTime int64  `yaml:"expireTime"`
			Cache      bool
		}
		Cluster struct {
			Token string
			List  []string
		}
		Version string
	}
	Log struct {
		Filename string
		Level    string
	}
}

var Configuration = configuration{}

func init() {
	conf := &configor.Config{
		AutoReload: true,
		AutoReloadCallback: func(config interface{}) {
			level, _ := logger.Parse(Configuration.Log.Level)
			logger.SetLevel(level)
		},
	}
	err := configor.New(conf).Load(&Configuration, "nacos.yml")
	if err != nil {
		logger.Panic(err.Error())
	}
}
