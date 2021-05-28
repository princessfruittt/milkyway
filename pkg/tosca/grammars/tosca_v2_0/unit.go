package tosca_v2_0

import (
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/normal"
)

//
// Unit
//
// See ServiceTemplate
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

type Unit struct {
	*Entity `name:"unit" yaml:",inline"`

	ToscaDefinitionsVersion *string           `read:"tosca_definitions_version" require:"" yaml:"tosca_definitions_version"`
	Profile                 *string           `read:"profile" yaml:"profile,omitempty"` // introduced in TOSCA 1.2 as "namespace", renamed in TOSCA 2.0
	Metadata                Metadata          `read:"metadata,!Metadata" yaml:"metadata,omitempty"`
	Description             *string           `read:"description" yaml:"description,omitempty"`
	Repositories            Repositories      `read:"repositories,Repository" yaml:"repositories,omitempty"`
	Imports                 Imports           `read:"imports,[]Import" yaml:"imports,omitempty"`
	ArtifactTypes           ArtifactTypes     `read:"artifact_types,ArtifactType" hierarchy:"" yaml:"artifact_types,omitempty"`
	CapabilityTypes         CapabilityTypes   `read:"capability_types,CapabilityType" hierarchy:"" yaml:"capability_types,omitempty"`
	DataTypes               DataTypes         `read:"data_types,DataType" hierarchy:"" yaml:"data_types,omitempty"`
	GroupTypes              GroupTypes        `read:"group_types,GroupType" hierarchy:"" yaml:"group_types,omitempty"`
	InterfaceTypes          InterfaceTypes    `read:"interface_types,InterfaceType" hierarchy:"" yaml:"interface_types,omitempty"`
	NodeTypes               MapNodeTypes      `read:"node_types,NodeType" hierarchy:"" yaml:"node_types,omitempty"`
	PolicyTypes             PolicyTypes       `read:"policy_types,PolicyType" hierarchy:"" yaml:"policy_types,omitempty"`
	RelationshipTypes       RelationshipTypes `read:"relationship_types,RelationshipType" hierarchy:"" yaml:"relationship_types,omitempty"`
}

func NewUnit(context *tosca.Context) *Unit {
	return &Unit{Entity: NewEntity(context)}
}

func (self Unit) AddArtifactType(aType ArtifactType) {
	self.ArtifactTypes = append(self.ArtifactTypes, &aType)
}

// tosca.Reader signature
func ReadUnit(context *tosca.Context) tosca.EntityPtr {
	self := NewUnit(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	if self.Profile != nil {
		context.CanonicalNamespace = self.Profile
	}
	return self
}

// parser.Importer interface
func (self *Unit) GetImportSpecs() []*tosca.ImportSpec {
	// TODO: importing should also import repositories

	var importSpecs = make([]*tosca.ImportSpec, 0, len(self.Imports))
	for _, import_ := range self.Imports {
		if importSpec, ok := import_.NewImportSpec(self); ok {
			importSpecs = append(importSpecs, importSpec)
		}
	}
	return importSpecs
}

func (self *Unit) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	logNormalize.Debug("unit")

	if self.Metadata != nil {
		for k, v := range self.Metadata {
			normalServiceTemplate.Metadata[k] = v
		}
	}
}
