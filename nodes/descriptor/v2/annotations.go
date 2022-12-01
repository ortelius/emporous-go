package v2

import (
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
)

type UpdateFunc func(node Node) error

// UpdateDescriptors updates descriptors and return updated descriptors with the modified
// v2 nodes.
func UpdateDescriptors(nodes []Node, updateFunc UpdateFunc) ([]ocispec.Descriptor, error) {
	var updateDescs []ocispec.Descriptor

	for _, node := range nodes {

		if err := updateFunc(node); err != nil {
			return nil, err
		}

		desc := node.Descriptor()

		mergedJSON, err := node.Properties.MarshalJSON()
		if err != nil {
			return nil, err
		}
		desc.Annotations[uorspec.AnnotationUORAttributes] = string(mergedJSON)

		updateDescs = append(updateDescs, desc)
	}
	return updateDescs, nil
}
