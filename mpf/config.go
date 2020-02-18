// Package mpf 配置
/**
 * 配置
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 9:22
 */
package mpf

import (
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/spf13/viper"
)

// ConfigViper viper配置
type ConfigViper struct {
    list       map[string]*viper.Viper
    dirConfigs string
}

// GetConfig 获取配置
func (config *ConfigViper) GetConfig(fileName string) *viper.Viper {
    conf, ok := config.list[fileName]
    if !ok {
        conf = viper.New()
        conf.AddConfigPath(config.dirConfigs)
        conf.SetConfigType("yaml")
        conf.SetConfigName(fileName)
        err := conf.ReadInConfig()
        if err != nil {
            panic(mperr.NewConfigViper(errorcode.ConfigViperGet, "获取"+fileName+"配置出错,"+err.Error(), err))
        }
        config.list[fileName] = conf
    }

    return conf
}

var (
    insConfig *ConfigViper
)

func init() {
    insConfig = &ConfigViper{make(map[string]*viper.Viper), ""}
}

// NewConfig 获取配置单例
func NewConfig() *ConfigViper {
    return insConfig
}
