package version

import (
	"fmt"
	"io"
	"runtime"
	"text/template"
)

var (
	// commit is the head commit from git
	commit string
	// buildDate in ISO8601 format
	buildDate string
	// version describes the version of the client
	// set at build time or detected during runtime.
	version = "v0.0.0-unknown"
	// buildData set at build time to add extra information
	// to the version.
	buildData string
)

var versionTemplate = `UOR Client:
 Version:	{{ .Version }}
 Go Version:	{{ .GoVersion }}
 Git Commit:	{{ .GitCommit }}
 Build Date:	{{ .BuildDate }}
 Platform:	{{ .Platform }}
`

type clientVersion struct {
	Platform  string
	Version   string
	GitCommit string
	GoVersion string
	BuildDate string
}

func GetVersion() string {
	return versionWithBuild()
}

// WriteVersion will output the templated version message.
func WriteVersion(writer io.Writer) error {
	versionInfo := clientVersion{
		Version:   versionWithBuild(),
		GitCommit: commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	tmp, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		return fmt.Errorf("template parsing error: %v", err)
	}

	return tmp.Execute(writer, versionInfo)
}

func versionWithBuild() string {
	if buildData != "" {
		return fmt.Sprintf("%s+%s", version, buildData)
	}

	return version
}
