package config

import (
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

type DiscordConfig struct {
	Token          string `mapstructure:"token"`
	DefaultChannel uint64 `mapstructure:"default_channel"`
}
type Config struct {
	Discord          DiscordConfig `mapstructure:"discord"`
	MessageSendToken string        `mapstructure:"message_send_token"`
}

var (
	configOnce   sync.Once
	globalConfig *Config
)

func InitConfig(configDir, configFile, configSuffix string) error {
	var err error
	err = nil
	configOnce.Do(func() {
		viper.SetConfigName(configFile)
		viper.SetConfigType(configSuffix)
		viper.AddConfigPath(configDir)

		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err = viper.ReadInConfig(); err != nil {
			err = fmt.Errorf("无法读取配置文件: %w", err)
		}
		logger.Infof("读取到的所有配置keys：%v", viper.AllKeys())
		globalConfig = &Config{}
		if err = viper.Unmarshal(&globalConfig); err != nil {
			err = fmt.Errorf("无法反序列化配置文件: %w", err)
		}
	})
	return err
}
func GetConfig() *Config {
	return globalConfig
}
