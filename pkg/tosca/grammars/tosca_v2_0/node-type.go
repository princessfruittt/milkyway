package tosca_v2_0

import (
	"milkyway/pkg/tosca"
)

//
// NodeType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.9
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.8
//

type NodeType struct {
	*Type `name:"node type" yaml:",inline"`

	PropertyDefinitions    PropertyDefinitions    `read:"properties,PropertyDefinition" inherit:"properties,Parent" yaml:"properties,omitempty"`
	AttributeDefinitions   AttributeDefinitions   `read:"attributes,AttributeDefinition" inherit:"attributes,Parent" yaml:"attributes,omitempty"`
	CapabilityDefinitions  CapabilityDefinitions  `read:"capabilities,CapabilityDefinition" inherit:"capabilities,Parent" yaml:"capabilities,omitempty"`
	RequirementDefinitions RequirementDefinitions `read:"requirements,{}RequirementDefinition" inherit:"requirements,Parent" yaml:"requirements,omitempty"` // sequenced list, but we read it into map
	InterfaceDefinitions   InterfaceDefinitions   `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent" yaml:"interfaces,omitempty"`
	ArtifactDefinitions    ArtifactDefinitions    `read:"artifacts,ArtifactDefinition" inherit:"artifacts,Parent" yaml:"artifacts,omitempty"`

	Parent *NodeType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
}

func NewNodeType(context *tosca.Context) *NodeType {
	return &NodeType{
		Type:                   NewType(context),
		PropertyDefinitions:    make(PropertyDefinitions),
		AttributeDefinitions:   make(AttributeDefinitions),
		CapabilityDefinitions:  make(CapabilityDefinitions),
		RequirementDefinitions: make(RequirementDefinitions),
		InterfaceDefinitions:   make(InterfaceDefinitions),
		ArtifactDefinitions:    make(ArtifactDefinitions),
	}
}

func (self NodeType) AddProperty(k string, v PropertyDefinition) {
	self.PropertyDefinitions[k] = &v
}
func (self NodeType) AddInterface(k string, v InterfaceDefinition) {
	self.InterfaceDefinitions[k] = &v
}
func (self NodeType) AddArtifact(k string, v ArtifactDefinition) {
	self.ArtifactDefinitions[k] = &v
}

// tosca.Reader signature
func ReadNodeType(context *tosca.Context) tosca.EntityPtr {
	self := NewNodeType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *NodeType) GetParent() tosca.EntityPtr {
	return self.Parent
}

// tosca.Inherits interface
func (self *NodeType) Inherit() {
	logInherit.Debugf("node type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(self.Parent.AttributeDefinitions)
	self.CapabilityDefinitions.Inherit(self.Parent.CapabilityDefinitions)
	self.RequirementDefinitions.Inherit(self.Parent.RequirementDefinitions)
	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)
	self.ArtifactDefinitions.Inherit(self.Parent.ArtifactDefinitions)
}

//
// NodeTypes
//

type NodeTypes []*NodeType
type MapNodeTypes map[string]*NodeType

func (self NodeTypes) IsCompatible(nodeType *NodeType) bool {
	for _, baseNodeType := range self {
		if baseNodeType.Context.Hierarchy.IsCompatible(baseNodeType, nodeType) {
			return true
		}
	}
	return false
}

func (self NodeTypes) ValidateSubset(subset NodeTypes, context *tosca.Context) bool {
	isSubset := true
	for _, subsetNodeType := range subset {
		if !self.IsCompatible(subsetNodeType) {
			context.ReportIncompatibleTypeInSet(subsetNodeType)
			isSubset = false
		}
	}
	return isSubset
}