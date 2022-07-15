# Schema

UOR Schema is the attribute type declaration for a UOR Collection. Schema also links application logic to a importing UOR Collection.

This document explains UOR Schema and the relationship of a UOR Collection with its imported Schema. 

## Schema Elements

There are three elements within a schema:

1. Attribute Type Declarations
2. Algorithm Reference
3. Default Attribute Mappings
4. Default Content Reference

### Attribute Type Declarations

Attribute type declarations MUST be written to the manifest annotations of the Schema Collection manifest. Attribute type declarations MUST follow the following syntax and guidelines:

**Key Formatting** - Attribute type declarations expressed via manifest annotation keys MUST follow reverse domain syntax, where the top level domain is the schema short name. 

**Values** - Attribute type declarations expressed via manifest annotation values MUST be one of the following:
  - a string
  - a number
  - a boolean
  - a dictionary 
  - null

### Algorithm Reference

Schema Collections MAY contain Algorithm References. The Algorithm Reference in a Schema Collection is the link to the algorithm imported into a calling Collection. This reference is expressed by assigning the `uor.algorithm=true` attribute to the node annotations of the Algorithm's Linked Collection. 

### Default Attribute Mappings

Schema Collections MAY contain Default Attribute Mappings. The Default Attribute Mappings in a schema instruct the UOR Client to add preset attributes to a Collection while being built. This reference is expressed by assigning the `uor.attribute.mapping=true` attribute to the node annotations of the Default Attribute Mappings. 

### Default Content Declaration

Schema Collections MAY contain a Default Content Declaration. The Default Content Declaration in a schema is referenced by an algorithm linked to a collection when the algorithm is rum. This Declaration is expressed by assigning the `uor.dcd={{ dictionary }}` attribute to the Manifest Annotations of the Schema Collection.  

## Design

Collections import Schema via an annotated Linked Collection. A Schema Collection imports an Algorithm into the Schema's calling collection. 

When the UOR Client retrieves a Collection, it first retrieves the OCI manifest of the referenced collection. The UOR Client then searches the manifest for a reference to an imported schema. If a schema is found, the UOR client retrieves the OCI Manifest of the imported schema. The UOR Client then searches the Schema Collection's OCI Manifest for an Algorithm Reference. If an Algorithm Reference is found, The UOR Client will first check its cache and if needed, download the Referenced Algorithm for further operations. 

A Collection can only import a single schema. However, a Collection may link to another collection with a different schema. There are two types of schema declarations in a Collection's OCI Manifest. Those are: `uor.schema={{ Schema Collection address (Full URI or just the digest of the referenced Schema Collection's OCI manifest) }}` and `uor.schema.linked={{ The digest of all Schema Collection OCI Manifest References inherited through Collection links}}`. When a collection links to another collection, all linked schemas of the Referenced Linked Collection are inherited by the linking Collection and written to the value of the `uor.schema.linked` attribute.



