package defaultmanager

import (
	"github.com/emporous/emporous-go/content"
	"github.com/emporous/emporous-go/log"
	"github.com/emporous/emporous-go/manager"
)

// DefaultManager is the default implementation for a collection manager.
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
