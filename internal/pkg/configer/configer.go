package configer

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

var Configer *Config

func readCfgFromFile(vp *viper.Viper, path string) error {
	if len(path) == 0 {
		return errors.New("config path is required")
	}
	vp.SetConfigFile(path)
	if err := vp.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func readCfgFromConsul(vp *viper.Viper, endpoint, path string) error {
	if err := vp.AddRemoteProvider("consul", endpoint, path); err != nil {
		return err
	}
	vp.SetConfigType("yaml")
	if err := vp.ReadRemoteConfig(); err != nil {
		return err
	}
	return nil
}

// LoadConfig 支持文件配置加载 其他方式 待完善
func LoadConfig(source, path string) (err error) {
	v := viper.New()
	switch source {
	case "consul":
		err = readCfgFromConsul(v, "", path)
	default:
		err = readCfgFromFile(v, path)
	}
	//to struct
	if err := v.Unmarshal(&Configer); err != nil {
		return err
	}
	log.Print(Configer)
	return nil
}
