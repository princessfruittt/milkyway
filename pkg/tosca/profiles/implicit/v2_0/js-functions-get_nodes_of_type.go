// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/functions/get_nodes_of_type.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.7.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.7.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.7.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.7.1

clout.exec('tosca.lib.utils');

function evaluate(typeName) {
	if (arguments.length !== 1)
		throw 'must have 1 argument';
	var names = [];
	for (var id in clout.vertexes) {
		var vertex = clout.vertexes[id];
		if (tosca.isTosca(vertex))
			names.push(vertex.properties.name);
	}
	return names;
}
`
}
