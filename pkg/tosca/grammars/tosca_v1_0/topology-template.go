package tosca_v1_0

import (
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/grammars/tosca_v2_0"
)

//
// TopologyTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.8
//

// tosca.Reader signature
func ReadTopologyTemplate(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("WorkflowDefinitions", "")

	return tosca_v2_0.ReadTopologyTemplate(context)
}
