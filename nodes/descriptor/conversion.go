package descriptor

import (
	empspec "github.com/emporous/collection-spec/specs-go/v1alpha1"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// CollectionToOCI converts collection descriptor to OCI descriptor.
func CollectionToOCI(desc empspec.Descriptor) (ocispec.Descriptor, error) {
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
func OCIToCollection(desc ocispec.Descriptor) (empspec.Descriptor, error) {
	if desc.Annotations == nil {
		desc.Annotations = map[string]string{}
	}

	attributes, err := AnnotationsToAttributes(desc.Annotations)
	if err != nil {
		return empspec.Descriptor{}, err
	}
	return empspec.Descriptor{
		MediaType:   desc.MediaType,
		Digest:      desc.Digest,
		Size:        desc.Size,
		URLs:        desc.URLs,
		Attributes:  attributes,
		Annotations: desc.Annotations,
	}, nil
}
