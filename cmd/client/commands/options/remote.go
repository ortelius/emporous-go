package options

import "github.com/spf13/pflag"

// Remote describes remote configuration options that can be set.
type Remote struct {
	Insecure  bool
	PlainHTTP bool
}

// BindFlags binds options from a flag set to Remote options.
func (o *Remote) BindFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.Insecure, "insecure", "", o.Insecure, "Allow connections to registries SSL registry without certs")
	fs.BoolVarP(&o.PlainHTTP, "plain-http", "", o.PlainHTTP, "Use plain http and not https when contacting registries")
}

// RemoteAuth describes remote authentication configuration options that can be set.
type RemoteAuth struct {
	Configs []string
}

// BindFlags binds options from a flag set to RemoteAuth options.
func (o *RemoteAuth) BindFlags(fs *pflag.FlagSet) {
	fs.StringArrayVarP(&o.Configs, "configs", "c", o.Configs, "Path(s) to your registry credentials. Defaults to well-known "+
		"auth locations ~/.docker/config.json and $XDG_RUNTIME_DIR/container/auth.json, in respective order.")
}
