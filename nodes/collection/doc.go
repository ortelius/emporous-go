/*
Copyright 2022 Emporous Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collection

// This package defines the Collection node type with loading and iteration strategies.

// The Collection type stores model.Node types in a graph structure.
// The reference implementation for the Collection node is to store the relationship
// between linked content within a collection of files (e.g. a collection of static files rendered by a web server).

// Basic iterator implementations are defines here within this package to iterate Collection child
// nodes as a slice in various orders.

// Since Collections can relate to other Collection instances, they implement the model.Node interface as well
// to allow storage and traversal of distributed, linked Collection nodes.
