package grammars

import (
	"milkyway/pkg/tosca"
)

// Map of keyword -> version -> grammar
var Grammars = make(map[string]map[string]*tosca.Grammar)

// Map of keyword -> version -> internal URL path
var ImplicitProfilePaths = make(map[string]map[string]string)
