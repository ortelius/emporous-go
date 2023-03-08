package components

import (
	"testing"

	"github.com/anchore/syft/syft/linux"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"

	"github.com/emporous/emporous-go/nodes/descriptor"
)

func TestInventoryToProperties(t *testing.T) {
	inv := makeSBOM()
	type spec struct {
		name       string
		inputProp  descriptor.Properties
		path       string
		assertFunc func(properties descriptor.Properties) bool
		expError   string
	}

	cases := []spec{
		{
			name:      "Success/EmptyProperties",
			inputProp: descriptor.Properties{},
			assertFunc: func(properties descriptor.Properties) bool {
				return properties.Descriptor != nil && properties.Descriptor.Name == "package-1"
			},
			path: "testpath-1",
		},
		{
			name: "Success/PropertiesMerge",
			inputProp: descriptor.Properties{
				Runtime: &ocispec.ImageConfig{
					User: "test",
				},
			},
			assertFunc: func(properties descriptor.Properties) bool {
				if properties.Descriptor == nil || properties.Descriptor.Name != "package-2" {
					return false
				}

				if properties.Runtime.User != "test" {
					return false
				}

				return true
			},
			path: "testpath-2",
		},
		{
			name:      "Success/PackageNotFound",
			inputProp: descriptor.Properties{},
			assertFunc: func(properties descriptor.Properties) bool {
				return properties.Descriptor == nil
			},
			path: "notthere",
		},
		{
			name:      "Failure/TooManyPackagesFound",
			inputProp: descriptor.Properties{},
			expError:  "incorrect number of components found for testpath-3, expected 1, got 2",
			path:      "testpath-3",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			props := &c.inputProp
			err := InventoryToProperties(inv, c.path, props)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.assertFunc(*props))
			}
		})
	}
}

func makeSBOM() sbom.SBOM {
	catalog := pkg.NewCatalog()
	location1 := source.NewLocation("testpath-1")
	catalog.Add(pkg.Package{
		Name:         "package-1",
		Version:      "1.0.1",
		Locations:    source.NewLocationSet(location1),
		Type:         pkg.PythonPkg,
		FoundBy:      "the-cataloger-1",
		Language:     pkg.Python,
		MetadataType: pkg.PythonPackageMetadataType,
		Licenses:     []string{"MIT"},
		Metadata: pkg.PythonPackageMetadata{
			Name:    "package-1",
			Version: "1.0.1",
		},
		PURL: "a-purl-1", // intentionally a bad pURL for test fixtures
		CPEs: []pkg.CPE{
			pkg.MustCPE("cpe:2.3:*:some:package:1:*:*:*:*:*:*:*"),
		},
	})
	location2 := source.NewLocation("testpath-2")
	catalog.Add(pkg.Package{
		Name:         "package-2",
		Version:      "2.0.1",
		Locations:    source.NewLocationSet(location2),
		Type:         pkg.DebPkg,
		FoundBy:      "the-cataloger-2",
		MetadataType: pkg.DpkgMetadataType,
		Metadata: pkg.DpkgMetadata{
			Package: "package-2",
			Version: "2.0.1",
		},
		PURL: "pkg:deb/debian/package-2@2.0.1",
		CPEs: []pkg.CPE{
			pkg.MustCPE("cpe:2.3:*:some:package:2:*:*:*:*:*:*:*"),
		},
	})
	location3 := source.NewLocation("testpath-3")
	catalog.Add(pkg.Package{
		Name:         "package-3",
		Version:      "3.0.1",
		Locations:    source.NewLocationSet(location3),
		Type:         pkg.DebPkg,
		FoundBy:      "the-cataloger-3",
		MetadataType: pkg.DpkgMetadataType,
		Metadata: pkg.DpkgMetadata{
			Package: "package-3",
			Version: "3.0.1",
		},
		PURL: "pkg:deb/debian/package-3@3.0.1",
		CPEs: []pkg.CPE{
			pkg.MustCPE("cpe:2.3:*:some:package:3:*:*:*:*:*:*:*"),
		},
	})
	catalog.Add(pkg.Package{
		Name:         "package-4",
		Version:      "4.0.1",
		Locations:    source.NewLocationSet(location3),
		Type:         pkg.DebPkg,
		FoundBy:      "the-cataloger-4",
		MetadataType: pkg.DpkgMetadataType,
		Metadata: pkg.DpkgMetadata{
			Package: "package-4",
			Version: "4.0.1",
		},
		PURL: "pkg:deb/debian/package-4@4.0.1",
		CPEs: []pkg.CPE{
			pkg.MustCPE("cpe:2.3:*:some:package:4:*:*:*:*:*:*:*"),
		},
	})
	return sbom.SBOM{
		Artifacts: sbom.Artifacts{
			PackageCatalog: catalog,
			LinuxDistribution: &linux.Release{
				PrettyName: "debian",
				Name:       "debian",
				ID:         "debian",
				IDLike:     []string{"like!"},
				Version:    "1.2.3",
				VersionID:  "1.2.3",
			},
		},
		Descriptor: sbom.Descriptor{
			Name:    "test",
			Version: "test",
		},
	}
}
