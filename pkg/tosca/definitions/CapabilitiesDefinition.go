// Copyright 2018 Bull S.A.S. Atos Technologies - Bull, Rue Jean Jaures, B.P.68, 78340, Les Clayes-sous-Bois, France.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package definitions

// An CapabilityDefinition is the representation of a TOSCA Capability Definition
// reference
// Simple YAML v1.2 http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#DEFN_ELEMENT_CAPABILITY_DEFN
// TOSCA 2.0 https://docs.oasis-open.org/tosca/TOSCA/v2.0/csd03/TOSCA-v2.0-csd03.html#_Toc56506365
// Type The mandatory name of the Capability Type this capability definition is based upon.
// Description The optional description of the Capability definition.
// Properties An optional map of property refinements for the Capability definition. The referred properties must have been defined in the Capability Type definition referred by the type keyword. New properties may not be add
// Attributes An optional map of attribute refinements for the Capability definition. The referred attributes must have been defined in the Capability Type definition referred by the type keyword. New attributes may not be added
// ValidSourceTypes An optional list of one or more valid names of Node Types that are supported as valid sources of any relationship established to the declared Capability Type. If undefined, all node types are valid sources.If valid_source_types is defined in the Capability Type, each element in this list must either be in or derived from an element in the list defined in the type
// Occurrences The optional minimum and maximum of occurrences for the capability. The occurrence represents the maximum number of relationships that are allowed by the Capability. If not defined the implied default is [1,UNBOUNDED] (which means that an exported Capability should allow at least one relationship to be formed with it and maximum a UNBOUNDED number of relationships).
type CapabilityDefinition struct {
	Type             string                         `yaml:"type" json:"type"`
	Description      string                         `yaml:"description,omitempty" json:"description,omitempty"`
	Properties       map[string]PropertyDefinition  `yaml:"properties,omitempty" json:"properties,omitempty"`
	Attributes       map[string]AttributeDefinition `yaml:"attributes,omitempty" json:"attributes,omitempty"`
	ValidSourceTypes []string                       `yaml:"valid_source_types,omitempty,flow" json:"valid_source_types,omitempty"`
	Occurrences      Range                          `yaml:"occurrences,omitempty" json:"occurrences,omitempty"`
}

//const (
//	// EndpointCapability is the default TOSCA type that should be used or
//	// extended to define a network endpoint capability.
//	EndpointCapability = "tosca.capabilities.Endpoint"
//
//	// PublicEndpointCapability represents a public endpoint.
//	PublicEndpointCapability = "tosca.capabilities.Endpoint.Public"
//)

//// An CapabilityAssignment is the representation of a TOSCA Capability Assignment
////
//// See http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.0/TOSCA-Simple-Profile-YAML-v1.0.html#DEFN_ELEMENT_CAPABILITY_ASSIGNMENT for more details
//type CapabilityAssignment struct {
//	Properties map[string]*ValueAssignment `yaml:"properties,omitempty" json:"properties,omitempty"`
//	Attributes map[string]*ValueAssignment `yaml:"attributes,omitempty" json:"attributes,omitempty"`
//}
//
//// UnmarshalYAML unmarshals a yaml into an CapabilityDefinition
//func (c *CapabilityDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
//	var s string
//	if err := unmarshal(&s); err == nil {
//		c.Type = s
//		return nil
//	}
//	var str struct {
//		Type             string                      `yaml:"type"`
//		Description      string                      `yaml:"description,omitempty"`
//		Properties       map[string]*ValueAssignment `yaml:"properties,omitempty"`
//		Attributes       map[string]*ValueAssignment `yaml:"attributes,omitempty"`
//		ValidSourceTypes []string                    `yaml:"valid_source_types,omitempty,flow"`
//		Occurrences      Range                       `yaml:"occurrences,omitempty"`
//	}
//	if err := unmarshal(&str); err != nil {
//		return err
//	}
//	c.Type = str.Type
//	c.Description = str.Description
//	c.Properties = str.Properties
//	c.Attributes = str.Attributes
//	c.ValidSourceTypes = str.ValidSourceTypes
//	c.Occurrences = str.Occurrences
//	return nil
//}
