package commands

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	managerapi "github.com/emporous/emporous-go/api/services/collectionmanager/v1alpha1"
	"github.com/emporous/emporous-go/cmd/client/commands/options"
	"github.com/emporous/emporous-go/content/layout"
	"github.com/emporous/emporous-go/manager/defaultmanager"
	"github.com/emporous/emporous-go/services/collectionmanager"
	"github.com/emporous/emporous-go/util/examples"
)

// ServeOptions describe configuration options that can
// be set using the serve subcommand.
type ServeOptions struct {
	*options.Common
	SocketLocation string
	options.Remote
}

var clientServeExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Serve with a specified unix domain socket location"},
	CommandString: "serve /var/run/test.sock",
}

// NewServeCmd creates a new cobra.Command for the serve subcommand.
func NewServeCmd(common *options.Common) *cobra.Command {
	o := ServeOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "serve SOCKET",
		Short:         "Serve gRPC API to allow Emporous collection management",
		Example:       examples.FormatExamples(clientServeExamples),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	o.Remote.BindFlags(cmd.Flags())

	return cmd
}

func (o *ServeOptions) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.SocketLocation = args[0]
	return nil
}

func (o *ServeOptions) Validate() error {
	return nil
}

func (o *ServeOptions) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interceptors := o.Logger.WithServerInterceptors()
	rpc := grpc.NewServer(grpc_middleware.WithUnaryServerChain(interceptors...))

	cache, err := layout.NewWithContext(ctx, o.CacheDir)
	if err != nil {
		return err
	}

	manager := defaultmanager.New(cache, o.Logger)

	opts := collectionmanager.ServiceOptions{
		Insecure:  o.Insecure,
		PlainHTTP: o.PlainHTTP,
		PullCache: cache,
	}
	service := collectionmanager.FromManager(manager, opts)

	// Register the service with the gRPC server
	managerapi.RegisterCollectionManagerServer(rpc, service)

	// Listen and serve
	lis, err := net.Listen("unix", o.SocketLocation)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-sigCh:
			s := <-sigCh
			o.Logger.Debugf("got signal %v, attempting graceful shutdown", s)
			cancel()
			rpc.GracefulStop()
		case <-ctx.Done():
		}
	}()

	if err := rpc.Serve(lis); err != nil {
		return err
	}

	wg.Wait()

	return nil

}
