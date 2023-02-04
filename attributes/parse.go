package attributes

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"

	"github.com/emporous/emporous-go/model"
)

// Parse will create a model.AttributeValue type from json.RawMessage.
func Parse(data json.RawMessage) (model.AttributeValue, error) {
	_, jsonType, _, err := jsonparser.Get(data)
	if err != nil {
		return nil, fmt.Errorf("somtring%w", err)
	}
	return getAttributeValue("", data, jsonType)
}

// getObject returns an object AttributeValue.
func getObject(data json.RawMessage) (model.AttributeValue, error) {
	set := map[string]model.AttributeValue{}
	handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) (err error) {
		keyAsString := string(key)
		attr, err := getAttributeValue(keyAsString, value, dataType)
		if err != nil {
			return err
		}
		set[keyAsString] = attr
		return nil
	}

	if err := jsonparser.ObjectEach(data, handler); err != nil {
		return nil, err
	}
	return NewObject(set), nil
}

// getList returns an list AttributeValue.
func getList(key string, data json.RawMessage) (model.AttributeValue, error) {
	var list []model.AttributeValue
	var pErr error

	handler := func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			pErr = err
			return
		}
		attr, err := getAttributeValue(key, value, dataType)
		if err != nil {
			pErr = err
			return
		}
		list = append(list, attr)
	}
	if pErr != nil {
		return nil, pErr
	}

	_, err := jsonparser.ArrayEach(data, handler)
	if err != nil {
		return nil, err
	}
	return NewList(list), nil
}

// getAttributeValue determines the json value type and return a corresponding AttributeValue.
func getAttributeValue(keyAsString string, value []byte, dataType jsonparser.ValueType) (model.AttributeValue, error) {
	valueAsString := string(value)
	var attr model.AttributeValue
	switch dataType {
	case jsonparser.String:
		trimmedVal := trimQuotes(valueAsString)
		attr = NewString(trimmedVal)
	case jsonparser.Number:
		// Using float for number like the standard lib
		floatVal, err := strconv.ParseFloat(valueAsString, 64)
		if err != nil {
			return attr, err
		}
		attr = NewFloat(floatVal)
	case jsonparser.Boolean:
		boolVal, err := strconv.ParseBool(valueAsString)
		if err != nil {
			return attr, err
		}
		attr = NewBool(boolVal)
	case jsonparser.Null:
		attr = NewNull()
	case jsonparser.Object:
		return getObject(value)
	case jsonparser.Array:
		return getList(keyAsString, value)
	default:
		return attr, ParseError{Key: keyAsString, Err: ErrInvalidAttribute}
	}
	return attr, nil
}

func trimQuotes(input string) string {
	return strings.Trim(input, `"`)
}

// ParseToSet converts a json.RawMessage to a model.AttributeSet.
func ParseToSet(input json.RawMessage) (model.AttributeSet, error) {
	set := map[string]model.AttributeValue{}

	values := map[string]json.RawMessage{}
	if err := json.Unmarshal(input, &values); err != nil {
		return nil, err
	}
	for key, val := range values {
		mattr, err := Parse(val)
		if err != nil {
			return nil, fmt.Errorf("error converting attribute %s to model: %v", key, err)
		}
		set[key] = mattr
	}
	return NewSet(set), nil
}
