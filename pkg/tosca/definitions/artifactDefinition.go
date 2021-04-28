// Copyright 2018 Bull S.A.S. Atos Technologies - Bull, Rue Jean Jaures, B.P.68, 78340, Les Clayes-sous-Bois, France.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package definitions

// An ArtifactDefinition is the representation of a TOSCA Artifact Definition
// An artifact definition defines a named, typed file that can be associated with Node Type or Node Template and used by orchestration engine to facilitate deployment and implementation of interface operations.
// Simple YAML v1.2 http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#DEFN_ENTITY_ARTIFACT_DEF for more details
// TOSCA 2.0 https://docs.oasis-open.org/tosca/TOSCA/v2.0/csd03/TOSCA-v2.0-csd03.html#_Toc56506476
// Type The mandatory artifact type for the artifact definition.
// File The mandatory URI string (relative or absolute) which can be used to locate the artifact’s file.
// Description The optional description for the artifact definition.
// Repository The optional name of the repository definition which contains the location of the external repository that contains the artifact.  The artifact is expected to be referenceable by its file URI within the repository.
// DeployPath The file path the associated file will be deployed on within the target node’s container.
// ArtifactVersion The version of this artifact. One use of this artifact_version is to declare the particular version of this artifact type, in addition to its mime_type (that is declared in the artifact type definition). Together with the mime_type it may be used to select a particular artifact processor for this artifact. For example, a python interpreter that can interpret python version 2.7.0.
// Checksum The checksum used to validate the integrity of the artifact.
// ChecksumAlgorithm Algorithm used to calculate the artifact checksum (e.g. MD5, SHA [Ref]). Shall be specified if checksum is specified for an artifact.
// Properties The optional map of property assignments associated with the artifact.
type ArtifactDefinition struct {
	Type              string                        `yaml:"type" json:"type"`
	File              string                        `yaml:"file" json:"file"`
	Description       string                        `yaml:"description,omitempty" json:"description,omitempty"`
	Repository        string                        `yaml:"repository,omitempty" json:"repository,omitempty"`
	DeployPath        string                        `yaml:"deploy_path,omitempty" json:"deploy_path,omitempty"`
	ArtifactVersion   string                        `yaml:"version,omitempty" json:"version,omitempty"`
	Checksum          string                        `yaml:"checksum,omitempty" json:"checksum,omitempty"`
	ChecksumAlgorithm string                        `yaml:"checksum_algorithm,omitempty" json:"checksum_algorithm,omitempty"`
	Properties        map[string]PropertyDefinition `yaml:"properties,omitempty" json:"properties,omitempty"`
}

//// ArtifactDefMap is a map of ArtifactDefinition
//type ArtifactDefMap map[string]ArtifactDefinition
//
//// UnmarshalYAML unmarshals a yaml into an ArtifactDefMap
//func (adm *ArtifactDefMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
//	log.Print("Resolving in artifacts in standard TOSCA format")
//	// Either a map or a seq
//	*adm = make(ArtifactDefMap)
//	var m map[string]ArtifactDefinition
//	if err := unmarshal(&m); err == nil {
//		for k, v := range m {
//			(*adm)[k] = v
//		}
//		return nil
//	}
//
//	log.Print("Resolving in artifacts in Alien format 1.2")
//	//var l []map[string]interface{}
//	var l []ArtifactDefinition
//	if err := unmarshal(&l); err == nil {
//
//		log.Print("list: %v", l)
//		for _, a := range l {
//			(*adm)[a.name] = a
//		}
//		return nil
//	}
//
//	log.Print("Resolving in artifacts in Alien format 1.3")
//	var lmap []ArtifactDefMap
//	if err := unmarshal(&lmap); err != nil {
//		return err
//	}
//	for _, m := range lmap {
//		for k, v := range m {
//			(*adm)[k] = v
//		}
//	}
//	return nil
//}
//
//
//
//// UnmarshalYAML unmarshals a yaml into an ArtifactDefinition
//func (a *ArtifactDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
//	var s string
//	if err := unmarshal(&s); err == nil {
//		a.File = s
//		return nil
//	}
//	var str struct {
//		Type        string `yaml:"type"`
//		File        string `yaml:"file"`
//		Description string `yaml:"description,omitempty"`
//		Repository  string `yaml:"repository,omitempty"`
//		DeployPath  string `yaml:"deploy_path,omitempty"`
//
//		// Extra types
//		MimeType string                 `yaml:"mime_type,omitempty"`
//		XXX      map[string]interface{} `yaml:",inline"`
//	}
//	if err := unmarshal(&str); err != nil {
//		return err
//	}
//	log.Print("Unmarshalled complex ArtifactDefinition %+v", str)
//	a.Type = str.Type
//	a.File = str.File
//	a.Description = str.Description
//	a.Repository = str.Repository
//	a.DeployPath = str.DeployPath
//	if str.File == "" && len(str.XXX) == 1 {
//		for k, v := range str.XXX {
//			a.name = k
//			var ok bool
//			a.File, ok = v.(string)
//			if !ok {
//				return errors.New("Missing mandatory attribute \"file\" for artifact")
//			}
//		}
//	}
//
//	if a.File == "" {
//		return errors.New("Missing mandatory attribute \"file\" for artifact")
//	}
//	return nil
//}
