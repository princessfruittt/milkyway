package cmd

import (
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"gopkg.in/yaml.v3"
	"log"
	tosca "milkyway/pkg/tosca"
	tosca2 "milkyway/pkg/tosca/grammars/tosca_v2_0"
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

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate [[Ansible role PATH or Github URL]] [flags]",
		Short: "generates TOSCA node type",
		Long:  "Generates valid tosca.node.type from Ansible role. ",
		Run:   cc.generateType,
	}
	cc.baseBuilderCmd = b.newBuilderCmd(cmd)
	return cc
}
func GitHubConnect(path string) (connection GithubConnection) {
	connection = *NewConnectionBuilder(path)
	err := connection.getContents("", "", "")
	if err != nil {
		log.Fatal(err)
	} else if NilFields(connection.ansibleRole) {
		log.Fatal(&cmdError{
			s:         "Please, make sure that Ansible role is correct",
			userError: true,
		})
	}
	return
}
func (c *generateCmd) generateType(cmd *cobra.Command, args []string) {
	if len(c.rolePath) == 0 && len(c.roleURL) == 0 {
		log.Fatal("Please, fill rolePath or roleURL flag ")
	} else {
		var generatedType *tosca2.NodeType
		if len(c.roleURL) > 0 {
			u, err := url.Parse(c.roleURL)
			if err != nil || len(u.Path) < 3 {
				log.Fatal(err)
			}
			con := GitHubConnect(u.Path)
			generatedType = con.ansibleRole.transform()
		} else if len(c.rolePath) > 0 {
			ansibleRole := GetRoleFromPath(c.rolePath)
			generatedType = ansibleRole.transform()
		}
		b, _ := yaml.Marshal(generatedType)
		log.Println(string(b))
	}
}

func (ar AnsibleRole) transform() *tosca2.NodeType {
	m := &ansibleRoleMeta{}
	errYaml := yaml.Unmarshal(ar.MetaMain, m)

	var someStruct map[string]interface{}
	errYaml2 := yaml.Unmarshal(ar.DefaultsMain, &someStruct)
	if errYaml != nil || errYaml2 != nil {
		log.Fatal(errYaml)
	}
	stylist := terminal.Stylize
	//if problemsFormat != "" {
	//	stylist = terminal.NewStylist(false)
	//}
	templateContext := tosca.NewContext(stylist, tosca.NewQuirks())
	newNodeType := tosca2.NewNodeType(templateContext)
	softwareComponent := tosca2.NewNodeType(templateContext)
	softwareComponent.Name = "tosca.nodes.SoftwareComponent"
	newNodeType.Type.Version = &tosca2.Version{
		CanonicalString: "",
		OriginalString:  "",
		Comparer:        "",
		Major:           0,
		Minor:           0,
		Fix:             0,
		Qualifier:       "",
		Build:           0,
	}
	newNodeType.Name = "ansible.nodes." + strcase.ToCamel(m.Meta.RoleName)
	newNodeType.Parent = softwareComponent
	newNodeType.Description = &m.Meta.Description
	newNodeType.Metadata = map[string]string{"author": m.Meta.Author, "min_ansible_version": m.Meta.MinAnsibleVersion}
	newNodeType.AddProperty("role_name", tosca2.PropertyDefinition{
		Required: &[]bool{true}[0], AttributeDefinition: &tosca2.AttributeDefinition{Name: "RoleName", DataTypeName: sPtr("string")},
		//Default: m.Meta.RoleName},
	})
	//for key, value := range someStruct {
	//	newNodeType.AddProperty(key,
	//		tosca2.PropertyDefinition{
	//			Required: &[]bool{true}[0],
	//			AttributeDefinition: &tosca2.AttributeDefinition{ DataTypeName: sPtr("string")}})
	//Default: value}})
	return newNodeType

}
func sPtr(s string) *string { return &s }
func NilFields(x AnsibleRole) bool {
	rv := reflect.ValueOf(&x).Elem()

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			return false
		}
	}
	return true
}
