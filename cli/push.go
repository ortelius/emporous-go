package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/client/registryclient/orasclient"
	"github.com/uor-framework/client/util/workspace"
)

type PushOptions struct {
	*RootOptions
	Destination string
	RootDir     string
	Insecure    bool
	PlainHTTP   bool
	Configs     []string
}

var clientPushExamples = templates.Examples(
	`
	# Push artifacts
	client push my-workspace localhost:5000/myartifacts:latest
	`,
)

func NewPushCmd(rootOpts *RootOptions) *cobra.Command {
	o := PushOptions{RootOptions: rootOpts}

	cmd := &cobra.Command{
		Use:           "push SRC DST",
		Short:         "Push OCI artifacts from specified source into a registry",
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
		o.Destination,
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

	descs, err := client.GatherDescriptors("", files...)
	if err != nil {
		return err
	}

	configDesc, err := client.GenerateConfig(nil)
	if err != nil {
		return err
	}

	if _, err := client.GenerateManifest(configDesc, nil, descs...); err != nil {
		return err
	}

	desc, err := client.Execute(ctx)
	if err != nil {
		return fmt.Errorf("error publishing content to %s: %v", o.Destination, err)
	}

	o.Logger.Infof("Artifact %s published to %s\n", desc.Digest, o.Destination)

	return nil
}
