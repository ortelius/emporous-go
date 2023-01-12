package descriptor

import (
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
)

// CollectionToOCI converts collection descriptor to OCI descriptor.
func CollectionToOCI(desc uorspec.Descriptor) (ocispec.Descriptor, error) {
	mergedAnnotations := map[string]string{}
	for key, value := range desc.Annotations {
		mergedAnnotations[key] = value
	}
	annotations, err := AnnotationsFromAttributes(desc.Attributes)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	for key, value := range annotations {
		if _, exist := mergedAnnotations[key]; !exist {
			mergedAnnotations[key] = value
		}
	}

	return ocispec.Descriptor{
		MediaType:   desc.MediaType,
		Digest:      desc.Digest,
		Size:        desc.Size,
		URLs:        desc.URLs,
		Annotations: desc.Annotations,
	}, nil
}

// OCIToCollection converts OCI descriptor to collection descriptor.
func OCIToCollection(desc ocispec.Descriptor) (uorspec.Descriptor, error) {
	if desc.Annotations == nil {
		desc.Annotations = map[string]string{}
	}

	attributes, err := AnnotationsToAttributes(desc.Annotations)
	if err != nil {
		return uorspec.Descriptor{}, err
	}
	return uorspec.Descriptor{
		MediaType:   desc.MediaType,
		Digest:      desc.Digest,
		Size:        desc.Size,
		URLs:        desc.URLs,
		Attributes:  attributes,
		Annotations: desc.Annotations,
	}, nil
}
