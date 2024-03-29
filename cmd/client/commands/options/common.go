package options

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/emporous/emporous-go/log"
)

// EnvConfig stores CLI runtime configuration from environment variables.
// Struct field names should match the name of the environment variable that the field is derived from.
type EnvConfig struct {
	EMPOROUS_DEV_MODE bool // true: show unimplemented stubs in --help
}

// Common describes global configuration options that can be set.
type Common struct {
	IOStreams genericclioptions.IOStreams
	LogLevel  string
	Logger    log.LoggerWithInterceptor
	CacheDir  string
	EnvConfig
}

// BindFlags binds options from a flag set to Common options.
func (o *Common) BindFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.LogLevel, "loglevel", "l", "info",
		"Log level (debug, info, warn, error, fatal)")
}

// Init initializes default values for Common options.
func (o *Common) Init() error {
	logger, err := log.NewLogrusLogger(o.IOStreams.Out, o.LogLevel)
	if err != nil {
		return err
	}
	o.Logger = logger

	cacheEnv := os.Getenv("EMPOROUS_CACHE")
	switch {
	case cacheEnv != "":
		o.CacheDir = cacheEnv
	default:
		o.CacheDir = filepath.Join(xdg.CacheHome, "emporous")
	}

	return nil
}
