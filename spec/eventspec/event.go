package eventspec

import (
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/nodes/collection"
)

type Event struct {
	ID string
}

type EventManager interface {
	SpawnEvent() error
	StopEvent(eventID string, timeout int64) error
	RemoveEvent(eventID string) error
	ListEvents() ([]Event, error)
	EventStatus(eventID string) error
	EventRequest() bn
}

type LinkedCollectionManager interface {
}

type CollectionManager interface {
	LinkedCollectionManager
	ListCollections(filter model.Matcher) ([]collection.Collection, error)
	CollectionStatus(collection.Collection)
	PullCollection() (string, error)
	RemoveCollection(collection.Collection) error
}

type User struct {
	collection.Collection
}

type ControlMessage struct {
	model.Attributes
	Digest      string
	Source      string
	Destination string
	Payload     bool
}

type ObjectRouter interface {
	NewEventRequest()
	SendToRouter()
	GetEngine()
	EngineInit()
	ObjectRetrieval()
	ObjectTransmit()
}

// NewObjectRequest sends a request for an
// object from the Router
func (c ControlMessage) NewObjectRequest() {

}

// EngineInit initiates communication with an event engine
func (c ControlMessage) EngineInit() {
	// Get default attributes from Collection's Schema

	// Lookup the object(s) with the default attributes.

	// Send a control message to the target event engine:
	// Set the control message attributes to the default
	// attributes of the schema and add the caller's attributes.
	// Set the control message source to the caller's address
	// Set the destination to the target event engine
	// The payload boolean is set to true if there is content
	// transmitted with the control message.
	// The digest is the hash of the control message attributes,
	// source address, and destination address together.

}

// NewEngineLookup performs an event engine lookup
// and then returns the URI of the event engine if found.
func (c ControlMessage) GetEngine() {
	// Use the Attributes from c to query the Collection located at the Destination from c.

	// Lazy pull the event engine root manifest and event engine

}

// NewEventRequest generates an event request control message.
func (c ControlMessage) NewEventRequest(a model.Attributes) (ControlMessage, error) {

	// Add callers attributes (a) to control message

	// Add default event engine address to control message

	// Add source address (attributes)

	// Add destination address (attributes)

	// Set Payload to false

	// The digest is set

	// Return Event Request Control Message
}

// SendToRouter sends a control message to the router
func (c ControlMessage) SendToRouter() error {

	// Send control message to router

}
