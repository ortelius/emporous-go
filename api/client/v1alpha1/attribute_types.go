package v1alpha1

// AttributeQueryKind object kind of AttributeQuery
const AttributeQueryKind = "AttributeQuery"

// AttributeQuery configures an attribute query against
// UOR collection content.
type AttributeQuery struct {
	TypeMeta `json:",inline"`
	// Attributes list the configuration for Attribute types.
	Attributes Attributes `json:"attributes"`
	LinkQuery  LinkQuery  `json:"links"`
	Digests    []string   `json:"digests"`
}

// LinkQuery configures a link query for a v3 compliant registry
type LinkQuery struct {
	// LinksTo finds are artifacts that link to this given digests
	LinksTo []string `json:"linksTo"`
	// LinksFrom finds and resolves (if applicable) all the links from a given digest
	LinksFrom []string `json:"linksFrom"`
	// FilterBy filters any found links results by attribute
	FilterBy Attributes `json:"filterBy"`
}

// Attributes is a map structure that holds all
// attribute information provided by the user.
type Attributes map[string]interface{}
