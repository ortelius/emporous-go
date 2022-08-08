package config

import (
	"github.com/spf13/viper"
	"path/filepath"

	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
)

// ReadCollectionConfig read the specified config into a CollectionConfiguration type.
func ReadCollectionConfig(configPath string) (v1alpha1.DataSetConfiguration, error) {
	var configuration v1alpha1.DataSetConfiguration
	cfg, err := readInConfig(configPath, configuration)
	if err != nil {
		return configuration, err
	}
	return cfg.(v1alpha1.DataSetConfiguration), nil
}

// ReadAttributeQuery read the specified config into a AttributeQuery type.
func ReadAttributeQuery(configPath string) (v1alpha1.AttributeQuery, error) {
	var configuration v1alpha1.AttributeQuery
	cfg, err := readInConfig(configPath, configuration)
	if err != nil {
		return configuration, err
	}
	return cfg.(v1alpha1.AttributeQuery), nil
}

func readInConfig(configPath string, object interface{}) (interface{}, error) {
	base := filepath.Base(configPath)
	dir := filepath.Dir(configPath)
	viper.SetConfigName(base)
	viper.AddConfigPath(filepath.Clean(dir))
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&object)
	if err != nil {
		return nil, err
	}
	return object, nil
}
