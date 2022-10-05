package registryclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindRegistry(t *testing.T) {
	type spec struct {
		name     string
		cfg      RegistryConfig
		inRef    string
		expError string
		expReg   Registry
	}
	cases := []spec{
		{
			name: "Success/OneMatch",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "*.example.com",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
					{
						Prefix: "*.not.com",
						Endpoint: Endpoint{
							SkipTLS: false,
						},
					},
				},
			},
			inRef: "reg.example.com",
			expReg: Registry{
				Prefix: "*.example.com",
				Endpoint: Endpoint{
					SkipTLS: true,
				},
			},
		},
		{
			name: "Success/MultipleMatches",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "*.example.com",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
					{
						Prefix: "*",
						Endpoint: Endpoint{
							SkipTLS: false,
						},
					},
				},
			},
			inRef: "reg.example.com",
			expReg: Registry{
				Prefix: "*.example.com",
				Endpoint: Endpoint{
					SkipTLS: true,
				},
			},
		},
		{
			name: "Success/SubDomainWildcard",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "reg.example.*",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
					{
						Prefix: "*",
						Endpoint: Endpoint{
							SkipTLS: false,
						},
					},
				},
			},
			inRef: "reg.example.com",
			expReg: Registry{
				Prefix: "reg.example.*",
				Endpoint: Endpoint{
					SkipTLS: true,
				},
			},
		},
		{
			name: "Success/NotMatch",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "*.not.com",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
				},
			},
			inRef:  "reg.example.com",
			expReg: Registry{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg, err := FindRegistry(c.cfg, c.inRef)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				if c.expReg.Prefix == "" {
					require.Equal(t, (*Registry)(nil), reg)
				} else {
					require.Equal(t, c.expReg, *reg)
				}
			}
		})
	}
}
