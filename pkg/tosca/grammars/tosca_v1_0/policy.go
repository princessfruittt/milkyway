package tosca_v1_0

import (
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/grammars/tosca_v2_0"
)

//
// Policy
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.6
//

// tosca.Reader signature
func ReadPolicy(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("Metadata", "")
	context.SetReadTag("TriggerDefinitions", "")

	return tosca_v2_0.ReadPolicy(context)
}
