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
	"strings"
)

type generateCmd struct {
	url  string
	path string
	name string
	*baseBuilderCmd
}

var importsList = []string{"relationships.yaml", "interfaces.yaml", "nodes.yaml", "artifacts.yaml"}
var pathMap = map[string]string{"tasks": "artifacts/tasks", "vars": "artifacts/vars", "handlers": "artifacts/handlers", "files": "artifacts/files", "defaults": "artifacts/defaults", "templates": "artifacts/templates"}

//TOSCA normative types
var scTypeName = "tosca.nodes.SoftwareComponent"
var aTypeName = "tosca.artifacts.Directory"
var iTypeName = "tosca.interfaces.node.lifecycle.Standard"

func sPtr(s string) *string { return &s }

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates TOSCA node type. Use with -- url or --path flag.",
		Long:  "Generates valid tosca.node.type from Ansible role. ",
		Run:   cc.generateType,
	}
	cc.baseBuilderCmd = b.newBuilderCmd(cmd)
	cc.cmd.PersistentFlags().StringVarP(&cc.url, "url", "u", "", "Ansible galaxy GitHub URL e.g. https://github.com/gantsign/ansible-role-golang")
	cc.cmd.PersistentFlags().StringVarP(&cc.path, "path", "p", "", "Path to the Directory with role e.g. myroles/ansible-role-golang")
	cc.cmd.PersistentFlags().StringVarP(&cc.name, "name", "n", "Default", "Ansible Role Name. Fill it, if role_name in meta/main.yaml is empty.")
	return cc
}
func (c *generateCmd) generateType(cmd *cobra.Command, args []string) {
	if len(c.path) == 0 && len(c.url) == 0 {
		log.Fatal("Please, fill -path or -url flag (i.e., generate -p [path_to_role]).")
	} else {
		var st *tosca2.ServiceTemplate
		if len(c.url) > 0 {
			u, err := url.Parse(c.url)
			if err != nil || len(u.Path) < 3 {
				log.Fatal(err)
			}
			con := GitHubConnect(u.Path)
			st = con.ansibleRole.transform(c.name)
		} else if len(c.path) > 0 {
			ansibleRole := GetRoleFromPath(c.path)
			st = ansibleRole.transform(c.name)
		}
		b, err := yaml.Marshal(st)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}

func (ar AnsibleRole) transform(name string) *tosca2.ServiceTemplate {
	ar.Name = strcase.ToCamel(name)
	stylist := terminal.Stylize
	//if problemsFormat != "" {
	//	stylist = terminal.NewStylist(false)
	//}
	templateContext := tosca.NewContext(stylist, tosca.NewQuirks())
	serviceTemplate := tosca2.NewServiceTemplate(templateContext)
	nodeType := tosca2.NewNodeType(templateContext)
	nodeType.Name = "ansible.nodes." + ar.Name
	if ar.MetaMain != nil {
		m := &ansibleRoleMeta{}
		err := yaml.Unmarshal(ar.MetaMain, m)
		if err != nil {
			log.Fatal(err)
		}
		var meta = make(map[string]string)
		if m.Meta.RoleName != "" {
			toCamel := strcase.ToCamel(m.Meta.RoleName)
			ar.Name = toCamel
			meta["role_name"] = toCamel
			nodeType.Name = "ansible.nodes." + toCamel
		}
		if m.Meta.Description != "" {
			nodeType.Description = &m.Meta.Description
		}

		if m.Meta.Author != "" {
			meta["author"] = m.Meta.Author
		}
		if m.Meta.MinAnsibleVersion != "" {
			meta["min_ansible_version"] = m.Meta.MinAnsibleVersion
		}
		nodeType.Metadata = meta
	}
	if ar.DefaultsMain != nil {
		var defaults map[string]interface{}
		err := yaml.Unmarshal(ar.DefaultsMain, &defaults)
		if err != nil {
			log.Fatal(err)
		}
		//fill properties
		for k, v := range defaults {
			nodeType.AddProperty(k,
				tosca2.PropertyDefinition{
					Required: new(bool),
					AttributeDefinition: &tosca2.AttributeDefinition{DataTypeName: sPtr("string"),
						DefaultString: v,
						Default:       &tosca2.Value{Entity: &tosca2.Entity{Context: &tosca.Context{Data: v}}}},
				})
		}
	}

	//nodeType.Type.Version = &tosca2.Version{
	//	CanonicalString: "",
	//	OriginalString:  "",
	//	Comparer:        "",
	//	Major:           0,
	//	Minor:           0,
	//	Fix:             0,
	//	Qualifier:       "",
	//	Build:           0,
	//}
	//&[]bool{true}[0]
	//fill meta
	nodeType.ParentName = &scTypeName
	//fill imports
	for _, s := range importsList {
		serviceTemplate.AddImport(&tosca2.Import{
			URL: sPtr(filepath.Join("normativetypes/2.0/", s))})
	}
	//fill requirements
	//fill interfaces
	if ar.TasksMain != nil {
		var dep []string = nil
		r := reflect.ValueOf(&ar).Elem()
		for k := range ar.Artifacts {
			if r.FieldByName(strings.Title(k)).IsValid() == true {
				dep = append(dep, k)
			}
			nodeType.AddArtifact(k, tosca2.ArtifactDefinition{File: sPtr(filepath.Join("artifacts", k)), ArtifactTypeName: &aTypeName})
		}
		nodeType.AddInterface("Standard", tosca2.InterfaceDefinition{
			InterfaceTypeName: &iTypeName,
			OperationDefinitions: map[string]*tosca2.OperationDefinition{"create": {
				Implementation: &tosca2.InterfaceImplementation{
					Primary:      sPtr(pathMap["tasks"] + "artifacts/tasks/main.yaml"),
					Dependencies: &dep,
				},
			}},
		})
	}
	serviceTemplate.Unit.NodeTypes = map[string]*tosca2.NodeType{}
	serviceTemplate.AddNodeType(nodeType.Name, *nodeType)
	serviceTemplate.AddDefinitionVersion()
	return serviceTemplate
}
