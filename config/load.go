package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"sigs.k8s.io/yaml"

	"github.com/emporous/emporous-go/api/client/v1alpha1"
)

// ReadDataSetConfig reads the specified config into a DataSetConfiguration type.
func ReadDataSetConfig(configPath string) (v1alpha1.DataSetConfiguration, error) {
	data, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return v1alpha1.DataSetConfiguration{}, err
	}

	return LoadDataSetConfig(data)
}

// LoadDataSetConfig loads a DataSetConfigurationType from input.
func LoadDataSetConfig(data []byte) (configuration v1alpha1.DataSetConfiguration, err error) {
	if data, err = yaml.YAMLToJSON(data); err != nil {
		return configuration, err
	}

	if err = checkMeta(data, v1alpha1.DataSetConfigurationKind); err != nil {
		return configuration, err
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&configuration); err != nil {
		return configuration, err
	}
	return configuration, err
}

// ReadSchemaConfig reads the specified config into a SchemaConfiguration type.
func ReadSchemaConfig(configPath string) (v1alpha1.SchemaConfiguration, error) {
	data, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return v1alpha1.SchemaConfiguration{}, err
	}

	return LoadSchemaConfig(data)
}

// LoadSchemaConfig loads a SchemaConfiguration type from input.
func LoadSchemaConfig(data []byte) (configuration v1alpha1.SchemaConfiguration, err error) {
	if data, err = yaml.YAMLToJSON(data); err != nil {
		return configuration, err
	}

	if err = checkMeta(data, v1alpha1.SchemaConfigurationKind); err != nil {
		return configuration, err
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&configuration); err != nil {
		return configuration, err
	}
	return configuration, err
}

// ReadAttributeQuery reads the specified config into a AttributeQuery type.
func ReadAttributeQuery(configPath string) (v1alpha1.AttributeQuery, error) {
	data, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return v1alpha1.AttributeQuery{}, err
	}

	return LoadAttributeQuery(data)
}

// LoadAttributeQuery loads an AttributeQuery type from input.
func LoadAttributeQuery(data []byte) (configuration v1alpha1.AttributeQuery, err error) {
	if data, err = yaml.YAMLToJSON(data); err != nil {
		return configuration, err
	}

	if err = checkMeta(data, v1alpha1.AttributeQueryKind); err != nil {
		return configuration, err
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&configuration); err != nil {
		return configuration, err
	}
	return configuration, err
}

func checkMeta(data []byte, kind string) error {
	var typeMeta v1alpha1.TypeMeta
	if err := json.Unmarshal(data, &typeMeta); err != nil {
		return fmt.Errorf("get type meta: %v", err)
	}
	if typeMeta.Kind != kind {
		return fmt.Errorf("config kind %s, does not match expected %s", typeMeta.Kind, kind)
	}
	return nil
}
