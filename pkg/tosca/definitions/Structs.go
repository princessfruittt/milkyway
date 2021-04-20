package definitions

import (
	"strings"
)

type Range struct {
	LowerBound uint64
	UpperBound uint64
}
type Input struct {
	ValueAssign *ValueAssignment    `json:"value_assignment,omitempty"`
	PropDef     *PropertyDefinition `json:"property_definition,omitempty"`
}
type Output struct {
	ValueAssign      *ValueAssignment               `json:"value_assignment,omitempty"`
	AttributeMapping map[string]AttributeDefinition `json:"attribute_mapping,omitempty"`
}

// An Implementation is the representation of the implementation part of a TOSCA Operation Definition
//
// See http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#DEFN_ELEMENT_OPERATION_DEF for more details
type Implementation struct {
	Primary       string             `yaml:"primary" json:"primary"`
	Dependencies  []string           `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Artifact      ArtifactDefinition `yaml:",inline" json:"artifact,omitempty"`
	OperationHost string             `yaml:"operation_host,omitempty" json:"operation_host,omitempty"`
}

// UNBOUNDED is the maximum value of a Range
// Max uint64 as per https://golang.org/ref/spec#Numeric_types
const UNBOUNDED uint64 = 18446744073709551615

func shouldQuoteYamlString(s string) bool {
	return strings.ContainsAny(s, ":[],\"{}#") ||
		strings.HasPrefix(strings.TrimSpace(s), "- ") ||
		strings.HasPrefix(strings.TrimSpace(s), "*") ||
		strings.HasPrefix(strings.TrimSpace(s), "?") ||
		strings.HasPrefix(strings.TrimSpace(s), "|") ||
		strings.HasPrefix(strings.TrimSpace(s), "!") ||
		strings.HasPrefix(strings.TrimSpace(s), "%") ||
		strings.HasPrefix(strings.TrimSpace(s), "@") ||
		strings.HasPrefix(strings.TrimSpace(s), "&")
}
