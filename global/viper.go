package global

import (
	"fmt"
	"github.com/spf13/viper"
)

func Viper() *viper.Viper {
	v := viper.New()
	var config = "/app/conf/app.yaml"
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fail to read yaml %s err:%v \n", config, err))
	}
	if err = v.Unmarshal(&CONF); err != nil {
		fmt.Println(err)
	}
	return v
}

