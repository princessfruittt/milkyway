package hot

import (
	"github.com/tliron/kutil/ard"
	"milkyway/pkg/tosca"
)

//
// Data
//

type Data struct {
	*Entity `name:"data"`

	Data ard.Value
}

func NewData(context *tosca.Context) *Data {
	return &Data{
		Entity: NewEntity(context),
		Data:   context.Data,
	}
}

// tosca.Reader signature
func ReadData(context *tosca.Context) tosca.EntityPtr {
	return NewData(context)
}
