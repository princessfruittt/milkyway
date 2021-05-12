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

// tosca.Contextual interface
func (self *Entity) GetContext() *tosca.Context {
	return self.Context
}
