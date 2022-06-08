package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/uor-framework/client/builder/api/v1alpha1"
)

func ReadConfig(configName string) (c v1alpha1.DataSetConfiguration, err error) {

	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	var configuration v1alpha1.DataSetConfiguration

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {

		fmt.Println("Config file not found")

	} else {
		err = viper.Unmarshal(&configuration)
		if err != nil {
			fmt.Printf("unable to decode into struct, %v", err)
		}
	}

	return configuration, err
}
