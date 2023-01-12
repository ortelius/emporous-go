package descriptor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/util/errlist"
)

var _ model.AttributeSet = &Properties{}

// Properties define all properties an UOR collection descriptor can have.
type Properties struct {
	Runtime    *ocispec.ImageConfig          `json:"core-runtime,omitempty"`
	Link       *uorspec.LinkAttributes       `json:"core-link,omitempty"`
	Descriptor *uorspec.DescriptorAttributes `json:"core-descriptor,omitempty"`
	Schema     *uorspec.SchemaAttributes     `json:"core-schema,omitempty"`
	File       *uorspec.File                 `json:"core-file,omitempty"`
	// A map of attribute sets where the string is the schema ID.
	Others map[string]model.AttributeSet `json:"-"`
}

// TODO(jpower432): When complex attribute sets are supported include core attributes filtering
// and searching in the AttributeSet methods.

// Exists checks for the existence of an attribute pair in the
// AttributeSet in the Properties.
// Only the "Others" field is evaluated during the search.
func (p *Properties) Exists(attribute model.Attribute) (bool, error) {
	for _, set := range p.Others {
		exists, err := set.Exists(attribute)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}

	return false, nil
}

// Find searches all AttributeSets in the Properties
// for a key and returns an attribute value.
// Only the "Others" field is evaluated during the search.
func (p *Properties) Find(s string) model.Attribute {
	for _, set := range p.Others {
		value := set.Find(s)
		if value != nil {
			return value
		}
	}
	return nil
}

// FindBySchema find the attribute value in the
// AttributeSet matching the given schema ID.
// Only the "Others" field is evaluated during the search.
func (p *Properties) FindBySchema(schema, key string) model.Attribute {
	set, found := p.Others[schema]
	if !found {
		return nil
	}
	return set.Find(key)
}

// ExistsBySchema checks the existence of the attributes in the
// AttributeSet matching the given schema ID.
// Only the "Others" field is evaluated during the search.
func (p *Properties) ExistsBySchema(schema string, attribute model.Attribute) (bool, error) {
	set, found := p.Others[schema]
	if !found {
		return false, nil
	}
	return set.Exists(attribute)
}

// MarshalJSON marshal an instance of Properties
// into the JSON format.
func (p *Properties) MarshalJSON() ([]byte, error) {
	propJSON, err := json.Marshal(*p)
	if err != nil {
		return nil, err
	}

	var mapping map[string]json.RawMessage
	if err = json.Unmarshal(propJSON, &mapping); err != nil {
		return nil, err
	}

	// Add attribute to the map without overriding struct fields
	for key, value := range p.Others {
		if _, ok := mapping[key]; ok {
			continue
		}
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		mapping[key] = valueJSON
	}

	return json.Marshal(mapping)
}

// List lists the AttributeSet attributes in the
// Properties. If the attribute under different schemas
// cannot merge, nil will be returned.
// Only the "Others" field is evaluated.
func (p *Properties) List() map[string]model.Attribute {
	var sets []model.AttributeSet
	for _, set := range p.Others {
		sets = append(sets, set)
	}
	mergedList, err := attributes.Merge(sets...)
	if err != nil {
		return nil
	}

	return mergedList.List()
}

// Len returns the length of the all AttributeSets
// in the Properties.
// Only the "Others" field is evaluated.
func (p *Properties) Len() int {
	var otherLen int
	for _, set := range p.Others {
		otherLen += set.Len()
	}
	return otherLen
}

// Merge merges a given AttributeSet into the descriptor Others AttributeSets.
func (p *Properties) Merge(sets map[string]model.AttributeSet) error {
	if len(sets) == 0 {
		return nil
	}

	for key, set := range sets {
		existingSet, exists := p.Others[key]
		if !exists {
			p.Others[key] = set
			continue
		}
		updatedSet, err := attributes.Merge(set, existingSet)
		if err != nil {
			return err
		}
		p.Others[key] = updatedSet
	}
	return nil
}

// IsALink returns whether a descriptor is a link.
func (p *Properties) IsALink() bool {
	return p.Link != nil
}

// IsASchema returns whether the descriptor is a schema.
func (p *Properties) IsASchema() bool {
	return p.Schema != nil
}

// IsAComponent returns whether the descriptor
// has a component name.
func (p *Properties) IsAComponent() bool {
	if p.Descriptor == nil {
		return false
	}
	return p.Descriptor.Component.Name != ""
}

// HasRuntimeInfo returns whether the descriptor
// has runtime information set.
func (p *Properties) HasRuntimeInfo() bool {
	return p.Runtime != nil
}

// HasFileInfo returns whether the descriptor
// has file information set.
func (p *Properties) HasFileInfo() bool {
	return p.File != nil
}

const (
	TypeLink       = "core-link"
	TypeDescriptor = "core-descriptor"
	TypeSchema     = "core-schema"
	TypeRuntime    = "core-runtime"
	TypeFile       = "core-file"
)

// Parse attempt to resolve attribute types in a set of json.RawMessage types
// into known types and adds unknown attributes to
// an attribute set, if supported.
func Parse(in map[string]json.RawMessage) (*Properties, error) {
	var out Properties
	other := map[string]model.AttributeSet{}

	var errs []error
	for key, prop := range in {
		switch key {
		case TypeLink:
			var l uorspec.LinkAttributes
			if err := json.Unmarshal(prop, &l); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
				continue
			}
			out.Link = &l
		case TypeDescriptor:
			var d uorspec.DescriptorAttributes
			if err := json.Unmarshal(prop, &d); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
				continue
			}
			out.Descriptor = &d
		case TypeSchema:
			var s uorspec.SchemaAttributes
			if err := json.Unmarshal(prop, &s); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
				continue
			}
			out.Schema = &s
		case TypeRuntime:
			var r ocispec.ImageConfig
			if err := json.Unmarshal(prop, &r); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
				continue
			}
			out.Runtime = &r
		case TypeFile:
			var f uorspec.File
			if err := json.Unmarshal(prop, &f); err != nil {
				errs = append(errs, ParseError{Key: key, Err: err})
				continue
			}
			out.File = &f
		default:
			set := attributes.Attributes{}
			handler := func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) (err error) {
				valueAsString := string(value)
				keyAsString := string(key)
				var attr model.Attribute
				switch dataType {
				case jsonparser.String:
					attr = attributes.NewString(keyAsString, valueAsString)
				case jsonparser.Number:
					// Using float for number like the standard lib
					floatVal, err := strconv.ParseFloat(valueAsString, 64)
					if err != nil {
						return err
					}
					attr = attributes.NewFloat(keyAsString, floatVal)
				case jsonparser.Boolean:
					boolVal, err := strconv.ParseBool(valueAsString)
					if err != nil {
						return err
					}
					attr = attributes.NewBool(keyAsString, boolVal)
				case jsonparser.Null:
					attr = attributes.NewNull(keyAsString)
				default:
					return ParseError{Key: keyAsString, Err: errors.New("unsupported attribute type")}
				}
				set[attr.Key()] = attr
				return nil
			}

			if err := jsonparser.ObjectEach(prop, handler); err != nil {
				errs = append(errs, fmt.Errorf("key %s: %w", key, err))
				continue
			}

			other[key] = set
		}
	}
	out.Others = other
	return &out, errlist.NewErrList(errs)
}
