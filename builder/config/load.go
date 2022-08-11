package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
	"io/ioutil"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

// ReadCollectionConfig read the specified config into a CollectionConfiguration type.
func ReadCollectionConfig(configPath string) (v1alpha1.DataSetConfiguration, error) {
	var configuration v1alpha1.DataSetConfiguration
	data, err := readInConfig(configPath, v1alpha1.DataSetConfigurationKind)
	if err != nil {
		return configuration, err
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&configuration); err != nil {
		return configuration, err
	}

	return configuration, nil
}

// ReadAttributeQuery read the specified config into a AttributeQuery type.
func ReadAttributeQuery(configPath string) (v1alpha1.AttributeQuery, error) {
	var configuration v1alpha1.AttributeQuery
	data, err := readInConfig(configPath, v1alpha1.AttributeQueryKind)
	if err != nil {
		return configuration, err
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&configuration); err != nil {
		return configuration, err
	}

	return configuration, nil
}

func readInConfig(configPath, kind string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return nil, err
	}

	if data, err = yaml.YAMLToJSON(data); err != nil {
		return nil, err
	}

	typeMeta, err := getTypeMeta(data)
	if err != nil {
		return nil, err
	}
	if typeMeta.Kind != kind {
		return nil, fmt.Errorf("config kind %s, does not match expected %s", typeMeta.Kind, kind)
	}
	return data, nil
}

func getTypeMeta(data []byte) (typeMeta v1alpha1.TypeMeta, err error) {
	if err := json.Unmarshal(data, &typeMeta); err != nil {
		return typeMeta, fmt.Errorf("get type meta: %v", err)
	}
	return typeMeta, nil
}
