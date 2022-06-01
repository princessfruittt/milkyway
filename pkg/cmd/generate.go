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
var pathMap = map[string]string{"tasks": "tasks.zip", "vars": "vars.zip", "handlers": "handlers.zip", "files": "files.zip", "defaults": "defaults.zip", "templates": "templates.zip"}

//TOSCA normative types
var scTypeName = "tosca.nodes.SoftwareComponent"
var cTypeName = "tosca.nodes.Compute"
var iTypeName = "tosca.interfaces.node.lifecycle.Standard"
var fTypeName = "tosca.artifacts.File"

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
		var ansibleRole AnsibleRole
		if len(c.url) > 0 {
			u, err := url.Parse(c.url)
			if err != nil || len(u.Path) < 3 {
				log.Fatal(err)
			}
			con := GitHubConnect(u.Path)
			ansibleRole = con.ansibleRole
		} else if len(c.path) > 0 {
			ansibleRole = GetRoleFromPath(c.path)
		}
		// ServiceTemplate similar for TOASCA versions
		var st *tosca2.ServiceTemplate
		st = ansibleRole.transform(c.name, c.toscaVersion)
		b, err := yaml.Marshal(st)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}

func (ar AnsibleRole) transform(name string, ver string) *tosca2.ServiceTemplate {
	ar.Name = strcase.ToCamel(name)
	stylist := terminal.Stylize
	templateContext := tosca.NewContext(stylist, tosca.NewQuirks())
	serviceTemplate := tosca2.NewServiceTemplate(templateContext)
	nodeType := tosca2.NewNodeType(templateContext)

	nodeType.Name = "ansible.nodes." + ar.Name
	nodeType.ParentName = &scTypeName

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
		if m.Meta.Platforms != nil && ver == "2.0" {
			nodeType = fillRequirements(nodeType, m.Meta.Platforms)
		}
	}

	if ar.DefaultsMain != nil {
		fillDefaultsProperties(nodeType, ar.DefaultsMain)
	}
	if ar.VarsMain != nil {
		fillVarsProperties(nodeType, ar.VarsMain)
	}
	if ar.Vars != nil {
		for _, v := range ar.Vars {
			fillVarsProperties(nodeType, v)
		}
	}
	if ar.TasksMain != nil {
		fillInterfaces(nodeType, ar)
	}
	switch ver {
	case "1.0", "1.1", "1.2", "1.3":
		serviceTemplate.AddDefinitionVersion("tosca_simple_yaml_" + strings.Replace(ver, ".", "_", -1))
	default:
		serviceTemplate.AddDefinitionVersion("tosca_2_0")
	}

	serviceTemplate.AddNodeType(nodeType.Name, *nodeType)
	//fillImports(serviceTemplate)
	return serviceTemplate
}

// func fillRequirements
func fillRequirements(nt *tosca2.NodeType, metaP []Platform) *tosca2.NodeType {
	platforms := map[string][]interface{}{}
	for _, p := range metaP {
		platforms[p.Name] = p.Versions
	}
	keys := make([]string, 0, len(platforms))
	values := make([][]interface{}, 0, len(platforms))

	for k, v := range platforms {
		keys = append(keys, k)
		values = append(values, v)
	}
	versions := map[string]string{}
	for _, v := range values {
		for _, e := range v {
			temp := fmt.Sprint(e)
			versions[temp] = temp
		}
	}
	vers := make([]string, 0, len(versions))
	for _, ver := range versions {
		vers = append(vers, ver)

	}
	cf := tosca2.CapabilityFilter{}
	cf = cf.AddPropertyFilters("distribution", tosca2.PropertyFilter{Valid: "[ " + strings.Join(keys, ", ") + " ]"})
	cf = cf.AddPropertyFilters("version", tosca2.PropertyFilter{Valid: "[ " + strings.Join(vers, ", ") + " ]"})
	nf := tosca2.NodeFilter{}
	nf = nf.AddCapabilityFilter("os", cf)
	req := tosca2.RequirementDefinition{
		TargetNodeTypeName: &cTypeName,
		TargetNodeFilter:   &nf,
	}
	nt = nt.AddRequirement("host", req)
	return nt
}

// func fillInterfaces check AnsibleRole folders with artifacts and fill tosca.Standard.create interface
func fillInterfaces(nt *tosca2.NodeType, ar AnsibleRole) {
	var dep []string = nil
	r := reflect.ValueOf(&ar).Elem()
	for k := range ar.Artifacts {
		if r.FieldByName(strings.Title(k)).IsValid() == true {
			dep = append(dep, k)
		}
		nt.AddArtifact(k, tosca2.ArtifactDefinition{File: sPtr(filepath.Join("artifacts", pathMap[k])), ArtifactTypeName: &fTypeName})
	}
	nt.AddInterface("Standard", tosca2.InterfaceDefinition{
		InterfaceTypeName: &iTypeName,
		OperationDefinitions: map[string]*tosca2.OperationDefinition{"create": {
			Implementation: &tosca2.InterfaceImplementation{
				Primary:      sPtr(filepath.Join("artifacts", "tasks", "main.yaml")),
				Dependencies: &dep,
			},
		}},
	})
}

// fillDefaultsProperties fill first entire of variables and add it to type properties.
func fillDefaultsProperties(nt *tosca2.NodeType, file []byte) {
	var vars map[string]interface{}
	err := yaml.Unmarshal(file, &vars)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range vars {
		nt.AddProperty(k,
			tosca2.PropertyDefinition{
				Required: new(bool),
				AttributeDefinition: &tosca2.AttributeDefinition{DataTypeName: sPtr("string"),
					DefaultString: v,
					Default:       &tosca2.Value{Entity: &tosca2.Entity{Context: &tosca.Context{Data: v}}}},
			})
	}
}

// func fillVarsProperties fill first entire of variables, and put repetitive variables to valid_value constraints
func fillVarsProperties(nt *tosca2.NodeType, file []byte) {
	var vars map[string]interface{}
	err := yaml.Unmarshal(file, &vars)
	if err != nil {
		log.Fatal(err)
	}
	var varsMap = make(map[string][]interface{})
	for p := range nt.PropertyDefinitions {
		for k, v := range vars {
			if p == k {
				varsMap[p] = append(varsMap[p], v)
			}
		}
		if varsMap[p] != nil {
			if len(nt.PropertyDefinitions[p].ConstraintClauses) > 0 {
				temp := strings.Join([]string{fmt.Sprint(varsMap[p]), nt.PropertyDefinitions[p].ConstraintClauses[0].Valid}, ",")
				nt.PropertyDefinitions[p].ConstraintClauses[0].Valid = temp
			} else {
				temp := []*tosca2.ConstraintClause{{Valid: fmt.Sprint(varsMap[p])}}
				nt.PropertyDefinitions[p].ConstraintClauses = temp
			}
		}
	}
	for k, v := range vars {
		if varsMap[k] == nil {
			nt.AddProperty(k,
				tosca2.PropertyDefinition{
					Required: new(bool),
					AttributeDefinition: &tosca2.AttributeDefinition{DataTypeName: sPtr("string"),
						DefaultString: v,
						Default:       &tosca2.Value{Entity: &tosca2.Entity{Context: &tosca.Context{Data: v}}}},
				})
		}
	}
}

// func fillArtifacts fill imports section for normative types import
func fillImports(st *tosca2.ServiceTemplate) {
	for _, s := range importsList {
		st.AddImport(&tosca2.Import{
			URL: sPtr(filepath.Join("normativetypes/2.0/", s))})
	}
}
