package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
	"reflect"
)

type generateCmd struct {
	*baseBuilderCmd
}
type AnsibleRole struct {
	TemplatesMain []byte
	TasksMain     []byte
	VarsMain      []byte
	DefaultsMain  []byte
	HandlersMain  []byte
	MetaMain      []byte
	FilesMain     []byte
	LibraryMain   []byte
	Templates     []byte
	Tasks         []byte
	Vars          []byte
	Defaults      []byte
	Handlers      []byte
	Meta          []byte
	Files         []byte
	Library       []byte
}
type ansibleRoleMeta struct {
	GalaxyInfo   GalaxyInfo `yaml:"galaxy_info,omitempty"`
	Dependencies []string   `yaml:"dependencies,omitempty"`
}
type GalaxyInfo struct {
	RoleName          string     `yaml:"role_name,omitempty"`
	Author            string     `yaml:"author,omitempty"`
	Description       string     `yaml:"description,omitempty"`
	Company           string     `yaml:"company,omitempty"`
	License           string     `yaml:"license,omitempty"`
	Platforms         []Platform `yaml:"platforms,omitempty"`
	MinAnsibleVersion string     `yaml:"min_ansible_version,omitempty"`
	GalaxyTags        []string   `yaml:"galaxy_tags,omitempty"`
}
type Platform struct {
	Name     string   `yaml:"name,omitempty"`
	Versions []string `yaml:"versions,omitempty"`
}

func NilFields(x AnsibleRole) bool {
	rv := reflect.ValueOf(&x).Elem()

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			return false
		}
	}
	return true
}

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate tosca types",
		Run:   cc.generateTypes,
	}
	cc.baseBuilderCmd = b.newBuilderCmd(cmd)
	return cc
}

func (c *generateCmd) generateTypes(cmd *cobra.Command, args []string) {
	u, err := url.Parse(c.roleURL)
	if err != nil || len(u.Path) < 3 {
		log.Fatal(err)
	}
	var connection = *NewConnectionBuilder(u.Path)

	//TODO Don't Check it
	ansibleRoleMeta := &ansibleRoleMeta{}
	//generatedType := tosca.NodeType{
	//	Type:         tosca.Type{},
	//	Properties:   nil,
	//	Attributes:   nil,
	//	Requirements: nil,
	//	Capabilities: nil,
	//	Interfaces:   nil,
	//	Artifacts:    nil,
	//}

	err = connection.getContents("", "")
	if err != nil {
		log.Fatal(err)
	} else if NilFields(connection.ansibleRole) {
		log.Fatal(&cmdError{
			s:         "Please, make sure that Ansible Role is correct",
			userError: true,
		})
	}

	errYaml := yaml.Unmarshal(connection.ansibleRole.MetaMain, ansibleRoleMeta)
	if errYaml != nil {
		log.Fatal(errYaml)
	} else {
		log.Print(ansibleRoleMeta)
		//log.Printf("--- m:\n%v\n\n", string(connection.ansibleRole.MetaMain))
	}
}
