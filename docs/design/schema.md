# Schema

UOR Schema is the attribute type declaration for a UOR Collection. Schema also links application logic to a importing UOR Collection. Schema can also be used for things like validating a Collection's links and attribute declarations.

This document explains UOR Schema and the relationship of a UOR Collection with its imported Schema. 

## Schema Elements

There are four elements within a schema:

1. Attribute Type Declarations
2. Algorithm Reference
3. Common Attribute Mappings
4. Default Content Reference

### Attribute Type Declarations

Attribute type declarations MUST reside within a node of the Schema Collection as a JSON Schema document. Attribute type
declarations MUST follow the following syntax and guidelines:

**Values** - Attribute type declarations are expression via manifest annotation. The values for limited to the following
and are expressions with JSON syntax:

1. Number (expressed at float)
2. Integer
3. String
4. Boolean
5. Null

### Algorithm Reference

Schemas Collections MAY contain Algorithm References. A Collection's Algorithm Reference can be thought of as the "
application logic" of the Collection. The Algorithm Reference in a Schema Collection is the link to the algorithm
imported into a calling Collection. This reference is expressed by assigning the `uor.algorithm=true` attribute to the
node annotations of the Algorithm's Linked Collection.

### Common Attribute Mappings

Schema Collections MAY contain Common Attribute Mappings. The Common Attribute Mappings in a schema instruct the UOR
Client to add preset attributes to a Collection while being built. This reference is expressed by assigning
the `uor.attribute.mapping=true` attribute to the node annotations of the Default Attribute Mappings.

### Default Content Declaration

Schema Collections MAY contain a Default Content Declaration. The Default Content Declaration in a schema is referenced by an algorithm linked to a collection when the algorithm is run. This Declaration is expressed by assigning the `uor.dcd={{ dictionary }}` attribute to the Manifest Annotations of the Schema Collection.  

## Design

Collections import Schema via an annotated Linked Collection. A Schema Collection imports an Algorithm into the Schema's calling collection. A Collection can only have one Schema and a Schema can have only one Algorithm Reference.  

When the UOR Client retrieves a Collection, it first retrieves the OCI manifest of the referenced collection. The UOR Client then searches the manifest for a reference to an imported schema. If a schema is found, the UOR client retrieves the OCI Manifest of the imported schema. The UOR Client then searches the Schema Collection's OCI Manifest for an Algorithm Reference. If an Algorithm Reference is found, The UOR Client will first check its cache and if needed, download the Referenced Algorithm for further operations. 

A Collection can only import a single schema. However, a Collection may link to another collection with a different schema. There are two types of schema declarations in a Collection's OCI Manifest. Those are: `uor.schema={{ Schema Collection address (Full URI or just the digest of the referenced Schema Collection's OCI manifest) }}` and `uor.schema.linked={{ The digest of all Schema Collection OCI Manifest References inherited through Collection links}}`. When a collection links to another collection, all linked schemas of the Referenced Linked Collection are inherited by the linking Collection and written to the value of the `uor.schema.linked` attribute.



