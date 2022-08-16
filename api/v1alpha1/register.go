package v1alpha1

import "path"

const (
	version = "v1alpha1"
	group   = "client.uor-framework.io"
)

var (
	// GroupVersion of UOR
	GroupVersion = path.Join(group, version)
)
