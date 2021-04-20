package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"milkyway/pkg/tosca"
	"milkyway/pkg/tosca/definitions"
	"net/url"
	"reflect"
)

type generateCmd struct {
	*baseBuilderCmd
}

const (
	ToscaString    = "string"
	ToscaInteger   = "integer"
	ToscaFloat     = "float"
	ToscaBool      = "boolean"
	ToscaTimestamp = "timestamp"
	ToscaNull      = "null"
)

// AnsibleRole struct for Ansible role v2.9 and less
type AnsibleRole struct {
	TemplatesMain []byte
	TasksMain     []byte
	VarsMain      []byte
	DefaultsMain  []byte
	HandlersMain  []byte
	MetaMain      []byte
	FilesMain     []byte
	//LibraryMain   []byte
	Templates []byte
	Tasks     []byte
	Vars      []byte
	Defaults  []byte
	Handlers  []byte
	Meta      []byte
	Files     []byte
	//Library       []byte
	Version string
}
type ansibleRoleMeta struct {
	Meta         GalaxyMeta `yaml:"galaxy_info,omitempty"`
	Dependencies []string   `yaml:"dependencies,omitempty"`
}

type GalaxyMeta struct {
	RoleName          string     `yaml:"role_name,omitempty"`
	Author            string     `yaml:"author,omitempty"`
	Description       string     `yaml:"description,omitempty"`
	Platforms         []Platform `yaml:"platforms,omitempty"`
	MinAnsibleVersion string     `yaml:"min_ansible_version,omitempty"`
}
type Platform struct {
	Name     string   `yaml:"name,omitempty"`
	Versions []string `yaml:"versions,omitempty"`
}

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate tosca types",
		Run:   cc.generateType,
	}
	cc.baseBuilderCmd = b.newBuilderCmd(cmd)
	return cc
}
func GitHubConnect(path string) (connection GithubConnection) {
	connection = *NewConnectionBuilder(path)
	err := connection.getContents("", "")
	if err != nil {
		log.Fatal(err)
	} else if NilFields(connection.ansibleRole) {
		log.Fatal(&cmdError{
			s:         "Please, make sure that Ansible Role is correct",
			userError: true,
		})
	}
	return
}
func (c *generateCmd) generateType(cmd *cobra.Command, args []string) {
	u, err := url.Parse(c.roleURL)
	if err != nil || len(u.Path) < 3 {
		log.Fatal(err)
	}
	//con := GitHubConnect(u.Path)
	//errYaml := yaml.Unmarshal(connection.ansibleRole.MetaMain, ansibleRoleMeta)

	var metaMain = `
dependencies: []

galaxy_info:
  role_name: nginx
  author: geerlingguy
  description: Nginx installation for Linux, FreeBSD and OpenBSD.
  company: "Midwestern Mac, LLC"
  license: "license (BSD, MIT)"
  min_ansible_version: 2.4
  platforms:
    - name: EL
      versions:
        - 7
        - 8
    - name: Debian
      versions:
        - all
    - name: Ubuntu
      versions:
        - trusty
        - xenial
        - focal
    - name: Archlinux
      versions:
        - all
    - name: FreeBSD
      versions:
        - 10.3
        - 10.2
        - 10.1
        - 10.0
        - 9.3
    - name: OpenBSD
      versions:
        - 5.9
        - 6.0
  galaxy_tags:
    - development
    - web
    - nginx
    - reverse
    - proxy
    - load
    - balancer
`

	//log.Print(ansibleRoleMeta)
	testrole := AnsibleRole{
		TemplatesMain: nil,
		TasksMain:     nil,
		VarsMain:      nil,
		DefaultsMain:  nil,
		HandlersMain:  nil,
		MetaMain:      nil,
		FilesMain:     nil,
		Templates:     nil,
		Tasks:         nil,
		Vars:          nil,
		Defaults:      nil,
		Handlers:      nil,
		Meta:          []byte(metaMain),
		Files:         nil,
		Version:       "",
	}
	generatedType := testrole.parseRole()
	//generatedType.Properties["role_name"].Default.Value = testrole.Meta.RoleName
	//m["ansible.role."+testrole.Meta.RoleName] = generatedType
	//log.Printf("--- m:\n%v\n\n", string(connection.ansibleRole.MetaMain))
	b, _ := yaml.Marshal(generatedType)
	log.Println(string(b))

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
func (ar AnsibleRole) parseRole() map[string]tosca.NodeType {
	m := &ansibleRoleMeta{}
	errYaml := yaml.Unmarshal([]byte(ar.Meta), m)
	if errYaml != nil {
		log.Fatal(errYaml)
	}

	nt := tosca.NodeType{
		Type: tosca.Type{
			Base: tosca.Node, DerivedFrom: "tosca.nodes.Root", Version: "0.0", Description: m.Meta.Description,
			Metadata: map[string]string{"author": m.Meta.Author, "min_ansible_version": m.Meta.MinAnsibleVersion}},
		Properties:   make(map[string]definitions.PropertyDefinition),
		Attributes:   nil,
		Requirements: nil,
		Capabilities: nil,
		Interfaces:   nil,
		Artifacts:    nil,
	}
	nt.AddProperty("role_name",
		definitions.PropertyDefinition{
			Type: ToscaString, Description: "RoleName", Required: &[]bool{true}[0],
			Default: &definitions.ValueAssignment{Type: definitions.ValueAssignmentMap, Value: m.Meta.RoleName}})
	return map[string]tosca.NodeType{m.Meta.RoleName: nt}
}
