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

// An AttributeDefinition is the representation of a TOSCA Attribute Definition
// An attribute definition defines a named, typed value that can be associated with an entity defined in this specification (e.g., a Node, Relationship or Capability Type).  Specifically, it is used to expose the “actual state” of some property of a TOSCA entity after it has been deployed and instantiated (as set by the TOSCA orchestrator). Attribute values can be retrieved via the get_attribute function from the instance model and used as values to other entities within TOSCA Service Templates.
// Simple YAML v1.2 http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#DEFN_ELEMENT_ATTRIBUTE_DEFN for more details
// TOSCA 2.0 https://docs.oasis-open.org/tosca/TOSCA/v2.0/csd03/TOSCA-v2.0-csd03.html#_Toc56506580
// Type The mandatory data type for the attribute.
// Description The optional description for the attribute.
// Default An optional key that may provide a value to be used as a default if not provided by another means.  This value SHALL be type compatible with the type declared by the attribute definition’s type keyname.
// Status The optional status of the attribute relative to the specification or implementation.  See supported status values . Defaults to supported.
// Constraints The optional list of sequenced constraint clauses for the attribute.
// KeySchema The schema definition for the keys used to identify entries in attributes of type TOSCA map (or types that derive from map). If not specified, the key_schema defaults to string. For attributes of type other than map, the key_schema is not allowed.
// EntrySchema The schema definition for the entries in attributes of TOSCA collection types such as list, map, or types that derive from list or map) If the attribute type is a collection type, the entry schema is mandatory. For other types, the entry_schema is not allowed.
// Metadata Defines a section used to declare additional metadata information.
type AttributeDefinition struct {
	Type        string           `yaml:"type" json:"type"`
	Description string           `yaml:"description,omitempty" json:"description,omitempty"`
	Default     *ValueAssignment `yaml:"default,omitempty" json:"default,omitempty"`
	Status      string           `yaml:"status,omitempty" json:"status,omitempty"`
	//Constraints []ConstraintClause     `yaml:"constraints,omitempty" json:"constraints,omitempty"`
	KeySchema   SchemaDefinition  `yaml:"key_schema,omitempty" json:"key_schema,omitempty"`
	EntrySchema SchemaDefinition  `yaml:"entry_schema,omitempty" json:"entry_schema,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

// UnmarshalYAML unmarshals a yaml into an AttributeDefinition
func (r *AttributeDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ra struct {
		Type        string           `yaml:"type"`
		Description string           `yaml:"description,omitempty"`
		Default     *ValueAssignment `yaml:"default,omitempty"`
		Status      string           `yaml:"status,omitempty"`
		EntrySchema SchemaDefinition `yaml:"entry_schema,omitempty"`
	}

	if err := unmarshal(&ra); err == nil && ra.Type != "" {
		r.Description = ra.Description
		r.Type = ra.Type
		r.Default = ra.Default
		r.Status = ra.Status
		r.EntrySchema = ra.EntrySchema
		return nil
	}

	var ras ValueAssignment
	if err := unmarshal(&ras); err != nil {
		return err
	}
	r.Type = ras.Type.String()
	r.Default = &ras
	return nil
}
