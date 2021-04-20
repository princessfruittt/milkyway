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

// An PropertyDefinition is the representation of a TOSCA Property Definition
// reference
// Simple YAML v1.2 http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#DEFN_ELEMENT_PROPERTY_DEFN for more details
// TOSCA 2.0 https://docs.oasis-open.org/tosca/TOSCA/v2.0/csd03/TOSCA-v2.0-csd03.html#BKM_Property_Def
type Status int

const (
	supported Status = iota
	unsupported
	experimental
	deprecated
)

type PropertyDefinition struct {
	Type        string      `yaml:"type" json:"type"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Required    *bool       `default:"true" yaml:"required,omitempty" json:"required,omitempty"`
	Default     interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	//Value
	Status Status `yaml:"status,omitempty" json:"status,omitempty"`
	//Constraints    []ConstraintClause `yaml:"constraints,omitempty"`
	KeySchema      SchemaDefinition  `yaml:"key_schema,omitempty" json:"key_schema,omitempty"`
	EntrySchema    SchemaDefinition  `yaml:"entry_schema,omitempty" json:"entry_schema,omitempty"`
	ExternalSchema string            `yaml:"external_schema,omitempty" json:"external_schema,omitempty"`
	Metadata       map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}
