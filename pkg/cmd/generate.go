package cmd

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"log"
	tosca "milkyway/pkg/tosca"
	tosca2 "milkyway/pkg/tosca/grammars/tosca_v2_0"
	"net/url"
	"path/filepath"
	"reflect"
)

type generateCmd struct {
	*baseBuilderCmd
}

var importsList = []string{"relationships.yaml", "interfaces.yaml", "nodes.yaml", "artifacts.yaml"}

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate",
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
		log.Fatal("Please, fill -rolePath or -roleURL flag (i.e., generate -p path_to_role).")
	} else {
		var generatedType *tosca2.ServiceTemplate
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
		b, err := yaml.Marshal(generatedType)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}

func (ar AnsibleRole) transform() *tosca2.ServiceTemplate {
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
	newTemplate := tosca2.NewServiceTemplate(templateContext)
	newNodeType := tosca2.NewNodeType(templateContext)
	//newNodeType.Type.Version = &tosca2.Version{
	//	CanonicalString: "",
	//	OriginalString:  "",
	//	Comparer:        "",
	//	Major:           0,
	//	Minor:           0,
	//	Fix:             0,
	//	Qualifier:       "",
	//	Build:           0,
	//}
	newNodeType.Name = "ansible.nodes." + strcase.ToCamel(m.Meta.RoleName)
	pn := "tosca.nodes.SoftwareComponent"
	newNodeType.ParentName = &pn
	newNodeType.Description = &m.Meta.Description
	newNodeType.Metadata = map[string]string{"author": m.Meta.Author, "min_ansible_version": m.Meta.MinAnsibleVersion, "role_name": m.Meta.RoleName}
	for k, v := range someStruct {
		tempA := tosca2.AttributeDefinition{DataTypeName: sPtr("string"),
			DefaultString: v,
			Default:       &tosca2.Value{Entity: &tosca2.Entity{Context: &tosca.Context{Data: v}}}}

		newNodeType.AddProperty(k,
			tosca2.PropertyDefinition{
				Required:            new(bool),
				AttributeDefinition: &tempA,
			})
	}
	//&[]bool{true}[0]
	for _, s := range importsList {
		tempS := filepath.Join("normativetypes/2.0/", s)
		newTemplate.AddImport(&tosca2.Import{
			URL: &tempS})
	}
	newTemplate.Unit.NodeTypes = map[string]*tosca2.NodeType{}
	newTemplate.AddNodeType(newNodeType.Name, *newNodeType)
	newTemplate.AddDefinitionVersion()
	return newTemplate

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
