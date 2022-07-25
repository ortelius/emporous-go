# Events

An Event is the entire interaction of the UOR Client with a referenced Collection. Since a Collection can have many nested Collections, a single event is the complete interaction with all Collections referenced by the UOR Client.

Collections have embedded Event Engines. An Event Engine is the application logic of a Collection. 

## Event Engine Lookup

1. UOR Client retrieves Collection root node
2. UOR Client looks up Schema in Collection root node
3. UOR Client retrieves the Schema Collection root node
4. UOR Client looks up application logic in the Schema Collection root node
5. UOR Client retrieves the application logic collection root node
6. UOR Client retrieves the application logic by platform/arch


## Router

When the UOR Client pushes or pulls (an object), it performs a routing action. The UOR Client can perform 2 routing actions:

1. Route object to Collection
2. Route object to Event Engine


### Router API
1. Route object to Collection Control Message
|-------------------|
| attributes (map)  | {The attributes of the object}
|-------------------|
| digest            | {The digest of the object}
|-------------------|
| source            | {The source address of the object}
|-------------------|
| destination       | {The destination address of the object}
|-------------------|
| payload (bool)    | {If no payload, the object is only attributes}
|-------------------|

# Event Engine Design

An Event Engine receives a Control Message and then performs the following tasks:

1. Prepends the source to all attributes in the Control Message
2. Incoming objects from the router are streamed into the needed format
3. Event Engine application is initiated

Note: It may be preferable to reference an object as a string formatted map (example: url encoded). This format can be used for things like naming object filenames when written to disk. This provides a convenient way for legacy applications to be ported to UOR. 









