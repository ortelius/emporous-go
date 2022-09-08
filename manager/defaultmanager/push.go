package defaultmanager

import (
	"context"
	"fmt"

	"github.com/uor-framework/uor-client-go/registryclient"
)

func (d DefaultManager) Push(ctx context.Context, reference string, remote registryclient.Remote) (string, error) {
	desc, err := remote.Push(ctx, d.store, reference)
	if err != nil {
		return "", fmt.Errorf("error publishing content to %s: %v", reference, err)
	}

	d.logger.Infof("Artifact %s published to %s\n", desc.Digest, reference)
	return desc.Digest.String(), nil
}
