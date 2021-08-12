package config

import (
	"frame/library/log"
	"github.com/spf13/viper"
)

type mysql struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Port int `yaml:"port"`
	Chart string `yaml:"chart"`
	Host string `yaml:"host"`
	DbName string `yaml:"db_name"`
}
type redis struct {
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Optins map[string]string `yaml:"optins"`
}
type configMap struct {
	Port int `yaml:"port"`
	Mode string `yaml:"mode"`
	LogPath string `yaml:"log_path"`
	Redis redis `yaml:"redis"`
	Mysql mysql `yaml:"mysql"`
}

var ConfigMap *configMap

// InitConfig 初始化配置
// Example (where serverCmd is a Cobra instance):
//
//	 config.InitConfig("./config.yaml")
//	 fmt.Println(config.ConfigMap.Mysql.Port)
//	 fmt.Println(viper.Get("mysql.port"))
// See (https://github.com/spf13/viper)
func InitConfig(filepath string)  {
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath(filepath)
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		log.LogOut.Error("load config err", struct {
			FilePath string
			ErrInfo string
		}{FilePath: filepath,ErrInfo: err.Error()})
	}
	if err := viper.Unmarshal(&ConfigMap); err != nil {
		log.LogOut.Error("Unmarshal config err", struct {
			FilePath string
			ErrInfo string
		}{FilePath: filepath,ErrInfo: err.Error()})
	}
}
