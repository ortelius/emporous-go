package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/client/builder/api/v1alpha1"
	load "github.com/uor-framework/client/builder/config"
	"github.com/uor-framework/client/registryclient"
	"github.com/uor-framework/client/registryclient/orasclient"
	"github.com/uor-framework/client/util/workspace"
)

// PushOptions describe configuration options that can
// be set using the push subcommand.
type PushOptions struct {
	*RootOptions
	Destination string
	RootDir     string
	Insecure    bool
	PlainHTTP   bool
	Configs     []string
	DSConfig    string
}

var clientPushExamples = templates.Examples(
	`
	# Push artifacts
	client push my-workspace localhost:5000/myartifacts:latest
	`,
)

// NewPushCmd creates a new cobra.Command for the push subcommand.
func NewPushCmd(rootOpts *RootOptions) *cobra.Command {
	o := PushOptions{RootOptions: rootOpts}

	cmd := &cobra.Command{
		Use:           "push SRC DST",
		Short:         "Push a UOR collection from specified source into a registry",
		Example:       clientPushExamples,
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringArrayVarP(&o.Configs, "configs", "c", o.Configs, "auth config paths")
	cmd.Flags().BoolVarP(&o.Insecure, "insecure", "", o.Insecure, "allow connections to SSL registry without certs")
	cmd.Flags().BoolVarP(&o.PlainHTTP, "plain-http", "", o.PlainHTTP, "use plain http and not https")
	cmd.Flags().StringVarP(&o.DSConfig, "dsconfig", "", o.DSConfig, "DataSet config path")

	return cmd
}

func (o *PushOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.RootDir = args[0]
	o.Destination = args[1]
	return nil
}

func (o *PushOptions) Validate() error {
	if _, err := os.Stat(o.RootDir); err != nil {
		return fmt.Errorf("workspace directory %q: %v", o.RootDir, err)
	}

	// TODO(jpower432): validate the reference and auth
	return nil
}

func (o *PushOptions) Run(ctx context.Context) error {
	space, err := workspace.NewLocalWorkspace(o.RootDir)
	if err != nil {
		return err
	}

	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithPlainHTTP(o.PlainHTTP),
		orasclient.WithAuthConfigs(o.Configs),
	)
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}

	var files []string
	err = space.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traversing %s: %v", path, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		if info.Mode().IsRegular() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	var config v1alpha1.DataSetConfiguration
	if len(o.DSConfig) > 0 {
		config, err = load.ReadConfig(o.DSConfig)
		if err != nil {
			return err
		}
	}

	// To allow the files to be loaded relative to the render
	// workspace, change to the render directory. This is required
	// to get path correct in the description annotations.
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(space.Path()); err != nil {
		return err
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			o.Logger.Errorf("%v", err)
		}
	}()

	descs, err := client.GatherDescriptors(ctx, "", files...)
	if err != nil {
		return err
	}

	// Add the attributes from the config to their respective blocks
	descs, err = AddDescriptors(descs, config)
	if err != nil {
		return err
	}

	configDesc, err := client.GenerateConfig(ctx, []byte("{}"), nil)
	if err != nil {
		return err
	}

	if _, err := client.GenerateManifest(ctx, o.Destination, configDesc, nil, descs...); err != nil {
		return err
	}

	desc, err := client.Execute(ctx, o.Destination, registryclient.TypePush)
	if err != nil {
		return fmt.Errorf("error publishing content to %s: %v", o.Destination, err)
	}

	o.Logger.Infof("Artifact %s published to %s\n", desc.Digest, o.Destination)

	return client.Destroy()
}

// AddDescriptors adds the attributes of each file listed in the config
// to the annotations of its respective descriptor.
func AddDescriptors(d []ocispec.Descriptor, c v1alpha1.DataSetConfiguration) ([]ocispec.Descriptor, error) {
	// For each descriptor
	for i1, desc := range d {
		// Get the filename of the block
		filename := desc.Annotations[ocispec.AnnotationTitle]
		// For each file in the config
		for i2, file := range c.Files {
			// If the config has a grouping declared, make a valid regex.
			if strings.Contains(file.File, "*") && !strings.Contains(file.File, ".*") {
				file.File = strings.Replace(file.File, "*", ".*", -1)
			} else {
				file.File = strings.Replace(file.File, file.File, "^"+file.File+"$", -1)
			}
			namesearch, err := regexp.Compile(file.File)
			if err != nil {
				return []ocispec.Descriptor{}, err
			}
			// Find the matching descriptor
			if namesearch.Match([]byte(filename)) {
				// Get the k/v pairs from the config and add them to the block's annotations.
				for k, v := range c.Files[i2].Attributes {
					d[i1].Annotations[k] = v
				}
			} else {
				// If the block does not have a corresponding config element, skip it.
				continue
			}
		}
	}
	return d, nil
}
