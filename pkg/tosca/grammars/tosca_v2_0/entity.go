package tosca_v2_0

import (
	"milkyway/pkg/tosca"
)

//
// Entity
//

type Entity struct {
	Context *tosca.Context `traverse:"ignore" json:"-" yaml:"-"`
}

func NewEntity(context *tosca.Context) *Entity {
	return &Entity{Context: context}
}
func NewEntityValue(context *tosca.Context, v interface{}) *Entity {
	ent := &Entity{Context: context}
	ent.Context.Data = v
	return ent
}

func (self *Entity) AddValue(v interface{}) {
	self.Context.Data = v
}

// tosca.Contextual interface
func (self *Entity) GetContext() *tosca.Context {
	return self.Context
}
