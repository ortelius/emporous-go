package v1alpha1

// TypeMeta contains type metadata.
// This structure is equivalent to k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta
type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}
