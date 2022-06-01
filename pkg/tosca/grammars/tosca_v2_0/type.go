package tosca_v2_0

import (
	"sync"

	"milkyway/pkg/tosca"
)

//
// Type
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.1
//

type Type struct {
	*Entity `json:"-" yaml:"-"`
	Name    string `namespace:"" yaml:"-"`

	ParentName  *string  `read:"derived_from" yaml:"derived_from"`
	Version     *Version `read:"version,version" yaml:"version,omitempty"`
	Metadata    Metadata `read:"metadata,!Metadata" yaml:"metadata,omitempty"`
	Description *string  `read:"description" yaml:"description,omitempty"`

	metadataLock sync.RWMutex
}

func NewType(context *tosca.Context) *Type {
	return &Type{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.HasMetadata interface
func (self *Type) GetDescription() (string, bool) {
	self.metadataLock.RLock()
	defer self.metadataLock.RUnlock()

	if self.Description != nil {
		return *self.Description, true
	}
	return "", false
}

// tosca.HasMetadata interface
func (self *Type) GetMetadata() (map[string]string, bool) {
	self.metadataLock.RLock()
	defer self.metadataLock.RUnlock()

	metadata := make(map[string]string)
	if self.Metadata != nil {
		for key, value := range self.Metadata {
			metadata[key] = value
		}
	}
	return metadata, true
}

// tosca.HasMetadata interface
func (self *Type) SetMetadata(name string, value string) bool {
	self.metadataLock.Lock()
	defer self.metadataLock.Unlock()

	self.Metadata[name] = value
	return true
}

func (self *Type) GetMetadataValue(key string) (string, bool) {
	if self.Metadata != nil {
		value, ok := self.Metadata[key]
		return value, ok
	}
	return "", false
}
