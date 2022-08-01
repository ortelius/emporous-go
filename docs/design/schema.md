# Schema

UOR Schema is the attribute type declaration for a UOR Collection. Schema also links application logic to a importing UOR Collection. Schema can also be used for things like validating a Collection's links and attribute declarations.

This document explains UOR Schema and the relationship of a UOR Collection with its imported Schema. 

## Schema Elements

There are four elements within a schema:

1. Attribute Type Declarations
2. Algorithm Reference
3. Default Attribute Mappings
4. Default Content Reference

### Attribute Type Declarations

Attribute type declarations MUST reside within a node of the Schema Collection as a JSON formatted document. Attribute type declarations MUST follow the following syntax and guidelines:

**Key Formatting** - Attributes expressed via manifest annotation keys MUST follow reverse domain notation, where the top level domain is the schema name; tag or digest. The Schema JSON does not declare the top level domain within. The top level domain is implicitly prepended as the top level domain by the UOR Client. 

**Example key name:** `quay.io/exampleOrg/exampleSchema:versionTag.category1.attribute1`

**Values** - Attribute type declarations expressed via manifest annotation values MUST be one of the following:
  - a string
  - a number
  - a boolean
  - a dictionary 
  - an array
  - null

### Algorithm Reference

Schema Collections MAY contain Event Engine References. A Collection's Event Engine can be thought of as the "application logic" of the Collection. The Event Engine reference in a Schema Collection is the link to the Event Engine imported into a calling Collection. This reference is expressed by assigning the `uor.event=true` attribute to the node annotations of the Event Engine's Linked Collection. 

### Default Attribute Mappings

Schema Collections MAY contain Default Attribute Mappings. The Default Attribute Mappings in a schema instruct the UOR Client to add preset attributes to a Collection while being built. This reference is expressed by assigning the `uor.attribute.mapping=true` attribute to the node annotations of the Default Attribute Mappings. 

### Default Content Declaration

Schema Collections MAY contain a Default Content Declaration. The Default Content Declaration in a schema is referenced by an algorithm linked to a collection when the algorithm is run. This Declaration is expressed by assigning the `uor.dcd={{ dictionary }}` attribute to the Manifest Annotations of the Schema Collection.  

## Design

Collections import Schema via an annotated Linked Collection. A Schema Collection imports an Event Engine into the Schema's calling collection. A Collection can only have one Schema and a Schema can have only one Event Engine Reference.  

When the UOR Client retrieves a Collection, it first retrieves the OCI manifest of the referenced collection. The UOR Client then searches the manifest for a reference to an imported schema. If a schema is found, the UOR client retrieves the OCI Manifest of the imported schema. The UOR Client then searches the Schema Collection's OCI Manifest for an Event Engine reference. If an Event Engine reference is found, The UOR Client will first check its cache and if needed, download the referenced Event Engine for further operations. 

A Collection can only import a single schema. However, a Collection may link to another collection with a different schema. There are two types of schema declarations in a Collection's OCI Manifest. Those are: `uor.schema={{ Schema Collection address (Full URI or just the digest of the referenced Schema Collection's OCI manifest) }}` and `uor.schema.linked={{ The digest of all Schema Collection OCI Manifest References inherited through Collection links}}`. When a collection links to another collection, all linked schemas of the Referenced Linked Collection are inherited by the linking Collection and written to the value of the `uor.schema.linked` attribute.



