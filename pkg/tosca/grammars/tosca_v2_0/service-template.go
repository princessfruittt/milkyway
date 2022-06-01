package tosca_v2_0

import (
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/normal"
)

//
// ServiceTemplate
//
// See Unit
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

type ServiceTemplate struct {
	*Unit `name:"service template" yaml:",inline"`

	TopologyTemplate *TopologyTemplate `read:"topology_template,TopologyTemplate" yaml:"topology_template,omitempty"`
}

func NewServiceTemplate(context *tosca.Context) *ServiceTemplate {
	return &ServiceTemplate{Unit: NewUnit(context)}
}

func (self ServiceTemplate) AddNodeType(k string, v NodeType) {
	self.Unit.NodeTypes[k] = &v
}

func (self ServiceTemplate) AddArtifactType(k string, v ArtifactType) {
	self.Unit.ArtifactTypes[k] = &v
}

func (self ServiceTemplate) AddImport(imp *Import) {
	temp := []*Import{imp}
	self.Unit.Imports = append(self.Unit.Imports, temp...)
}

func (self ServiceTemplate) AddDefinitionVersion(version string) {
	self.Unit.ToscaDefinitionsVersion = &version
}

func (self ServiceTemplate) AddUnit(u Unit) {
	self.Unit = &u
}

// tosca.Reader signature
func ReadServiceTemplate(context *tosca.Context) tosca.EntityPtr {
	self := NewServiceTemplate(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	if self.Profile != nil {
		context.CanonicalNamespace = self.Profile
	}
	return self
}

// normal.Normalizable interface
func (self *ServiceTemplate) NormalizeServiceTemplate() *normal.ServiceTemplate {
	logNormalize.Debug("service template")

	normalServiceTemplate := normal.NewServiceTemplate()

	if self.Description != nil {
		normalServiceTemplate.Description = *self.Description
	}

	normalServiceTemplate.ScriptletNamespace = self.Context.ScriptletNamespace

	self.Unit.Normalize(normalServiceTemplate)
	if self.TopologyTemplate != nil {
		self.TopologyTemplate.Normalize(normalServiceTemplate)
	}

	return normalServiceTemplate
}
