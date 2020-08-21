package cfg

import (
	"github.com/spf13/viper"
	"log"
)

func NewConfigParserDefault() *ConfParser {
	return &ConfParser{
		Path:             []string{DEFAULT_CONFIG_PATH, DEFAULT_CONFIG_PATH_MULTIAPP},
		ConfigName:       DEFAULT_CONFIG_FILE,
		OnConfigChangeCb: func() {},
	}
}

func (Self *ConfParser) AddPath(path string) {
	Self.Path = append(Self.Path, path)
}

func (Self *ConfParser) InitConf() {
	viper.SetConfigName(Self.ConfigName)
	for _, p := range Self.Path {
		viper.AddConfigPath(p)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件错误: %s\n", err.Error())
	}
	Self.viper = viper.GetViper()
}

func (Self *ConfParser) Reformat(in interface{}) error {
	return Self.viper.Unmarshal(&in)
}
