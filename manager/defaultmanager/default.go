package defaultmanager

import (
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/log"
	"github.com/uor-framework/uor-client-go/manager"
)

// DefaultManager is the default implementation for a collection router.
type DefaultManager struct {
	store  content.AttributeStore
	logger log.Logger
}

// New instantiates a new DefaultManager.
func New(store content.AttributeStore, logger log.Logger) manager.Manager {
	return DefaultManager{
		store:  store,
		logger: logger,
	}
}
