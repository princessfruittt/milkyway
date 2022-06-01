package tosca_v2_0

import (
	"github.com/tliron/kutil/util"
	"milkyway/pkg/tosca"
)

//
// Bytes
//
// [TOSCA-v2.0] @ ?
//

type Bytes struct {
	OriginalString string `json:"$originalString" yaml:"$originalString"`

	Bytes []byte `json:"bytes" yaml:"bytes"`
}

// tosca.Reader signature
func ReadBytes(context *tosca.Context) tosca.EntityPtr {
	var self Bytes

	if b64 := context.ReadString(); b64 != nil {
		self.OriginalString = *b64
		var err error
		if self.Bytes, err = util.FromBase64(self.OriginalString); err != nil {
			context.ReportValueMalformed("bytes", err.Error())
		}
	}

	return &self
}