package v1alpha1

// TypeMeta contains type metadata.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}
