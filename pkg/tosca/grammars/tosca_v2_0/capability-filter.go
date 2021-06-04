package tosca_v2_0

import (
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/normal"
)

//
// CapabilityFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.5.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.5.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.4.2
//

type CapabilityFilter struct {
	*Entity `name:"capability filter" yaml:"-"`
	Name    string `yaml:"-"`

	PropertyFilters PropertyFilters `read:"properties,PropertyFilter" yaml:",inline"`
}

func (self CapabilityFilter) AddPropertyFilters(k string, v PropertyFilter) CapabilityFilter {
	if self.PropertyFilters == nil {
		self.PropertyFilters = make(PropertyFilters)
	}
	self.PropertyFilters[k] = &v
	return self
}
func NewCapabilityFilter(context *tosca.Context) *CapabilityFilter {
	return &CapabilityFilter{
		Entity:          NewEntity(context),
		Name:            context.Name,
		PropertyFilters: make(PropertyFilters),
	}
}

// tosca.Reader signature
func ReadCapabilityFilter(context *tosca.Context) tosca.EntityPtr {
	self := NewCapabilityFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self CapabilityFilter) Normalize(normalRequirement *normal.Requirement) normal.FunctionCallMap {
	if len(self.PropertyFilters) == 0 {
		return nil
	}

	var normalFunctionCallMap normal.FunctionCallMap
	var ok bool
	if normalFunctionCallMap, ok = normalRequirement.CapabilityPropertyConstraints[self.Name]; !ok {
		normalFunctionCallMap = make(normal.FunctionCallMap)
		normalRequirement.CapabilityPropertyConstraints[self.Name] = normalFunctionCallMap
	}

	self.PropertyFilters.Normalize(normalFunctionCallMap)

	return normalFunctionCallMap
}

//
// CapabilityFilters
//

type CapabilityFilters []*CapabilityFilter

func (self CapabilityFilters) Normalize(normalRequirement *normal.Requirement) {
	for _, capabilityFilter := range self {
		capabilityFilter.Normalize(normalRequirement)
	}
}
