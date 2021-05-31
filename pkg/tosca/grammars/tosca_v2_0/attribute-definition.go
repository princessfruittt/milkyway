package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/normal"
)

//
// AttributeDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.10
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.10
//

type AttributeDefinition struct {
	*Entity  `name:"attribute definition" yaml:"-"`
	Name     string    `yaml:"-"`
	DataType *DataType `lookup:"type,DataTypeName" json:"-" yaml:"-"`

	DataTypeName      *string           `read:"type" yaml:"type,omitempty"` // required only if cannot be inherited or discovered
	Description       *string           `read:"description" yaml:"description,omitempty"`
	Default           *Value            `read:"default,Value" yaml:"-"`
	DefaultString     interface{}       `read:"default" yaml:"default"`
	Status            *string           `read:"status,omitempty" yaml:"status,omitempty"`
	ConstraintClauses ConstraintClauses `read:"constraints,[]ConstraintClause" yaml:"-"`
	KeySchema         *Schema           `read:"key_schema,Schema" yaml:"key_schema,omitempty"`     // introduced in TOSCA 1.3
	EntrySchema       *Schema           `read:"entry_schema,Schema" yaml:"entry_schema,omitempty"` // required if list or map
	Metadata          Metadata          `read:"metadata,Metadata" yaml:"metadata,omitempty"`       // introduced in TOSCA 1.3

	rendered bool
}

func NewAttributeDefinition(context *tosca.Context) *AttributeDefinition {
	return &AttributeDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadAttributeDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewAttributeDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *AttributeDefinition) GetKey() string {
	return self.Name
}

func (self *AttributeDefinition) Inherit(parentDefinition *AttributeDefinition) {
	logInherit.Debugf("attribute definition: %s", self.Name)

	// Validate type compatibility
	if (self.DataType != nil) && (parentDefinition.DataType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.DataType, self.DataType) {
		self.Context.ReportIncompatibleType(self.DataType, parentDefinition.DataType)
		return
	}

	if (self.Description == nil) && (parentDefinition.Description != nil) {
		self.Description = parentDefinition.Description
	}
	if (self.DataTypeName == nil) && (parentDefinition.DataTypeName != nil) {
		self.DataTypeName = parentDefinition.DataTypeName
	}
	if (self.KeySchema == nil) && (parentDefinition.KeySchema != nil) {
		self.KeySchema = parentDefinition.KeySchema
	}
	if (self.EntrySchema == nil) && (parentDefinition.EntrySchema != nil) {
		self.EntrySchema = parentDefinition.EntrySchema
	}
	if (self.Default == nil) && (parentDefinition.Default != nil) {
		self.Default = parentDefinition.Default
	}
	if (self.Status == nil) && (parentDefinition.Status != nil) {
		self.Status = parentDefinition.Status
	}
	if (self.DataType == nil) && (parentDefinition.DataType != nil) {
		self.DataType = parentDefinition.DataType
	}
}

// parser.Renderable interface
func (self *AttributeDefinition) Render() {
	logRender.Debugf("attribute definition: %s", self.Name)

	if self.rendered {
		// Avoid rendering more than once (can happen if we were called from Value.RenderAttribute)
		return
	}
	self.rendered = true

	if self.DataTypeName == nil {
		self.Context.FieldChild("type", nil).ReportFieldMissing()
		return
	}

	if self.DataType == nil {
		return
	}

	if internalTypeName, ok := self.DataType.GetInternalTypeName(); ok {
		switch internalTypeName {
		case ard.TypeList, ard.TypeMap:
			if self.EntrySchema == nil {
				self.EntrySchema = self.DataType.EntrySchema
			}

			// Make sure we have an entry schema
			if (self.EntrySchema == nil) || (self.EntrySchema.DataType == nil) {
				self.Context.ReportMissingEntrySchema(self.DataType.Name)
				return
			}

			if internalTypeName == ard.TypeMap {
				if self.KeySchema == nil {
					self.KeySchema = self.DataType.KeySchema
				}

				if self.KeySchema == nil {
					// Default to "string" for key schema
					self.KeySchema = ReadSchema(self.Context.FieldChild("key_schema", "string")).(*Schema)
					if !self.KeySchema.LookupDataType() {
						return
					}
				}
			}
		}
	}

	if self.Default != nil {
		// The "default" value must be a valid value of the type
		self.Default.RenderAttribute(self.DataType, self, false, false)
	}
}

func (self *AttributeDefinition) GetTypeInformation() *normal.TypeInformation {
	information := normal.NewTypeInformation()
	information.Metadata = tosca.GetInformationMetadata(self.Metadata)
	if self.Description != nil {
		information.Description = *self.Description
	}
	return information
}

//
// AttributeDefinitions
//

type AttributeDefinitions map[string]*AttributeDefinition

func (self AttributeDefinitions) Inherit(parentDefinitions AttributeDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}
