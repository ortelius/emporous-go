package v1alpha1

import "path"

const (
	version = "v1alpha1"
	group   = "client.emporous.io"
)

var (
	// GroupVersion of Emporous
	GroupVersion = path.Join(group, version)
)
