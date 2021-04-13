package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"net/url"
)

type generateCmd struct {
	*baseBuilderCmd
}

type ToscaNodeType struct {
	Properties string `yaml:"properties"`
}

type toscaTemplate struct {
	Version       string `yaml:"tosca_definitions_version"`
	ToscaNodeType `yaml:",inline"`
}
type ansibleRole struct {
	templates []byte
	tasks     []byte
	vars      []byte
	defaults  []byte
	handlers  []byte
	meta      []byte
	files     []byte
	library   []byte
}

var toscatemplate = toscaTemplate{
	Version: "tosca_simple_yaml_1_3",
	ToscaNodeType: ToscaNodeType{
		Properties: "params",
	},
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
	somedata := &toscaTemplate{}

	yaml_temp, _ := yaml.Marshal(toscatemplate)
	log.Println(yaml.Unmarshal(yaml_temp, somedata))
	errYaml := yaml.Unmarshal(yaml_temp, somedata)
	if errYaml != nil {
		log.Fatal(errYaml)
	} else {
		log.Print(somedata)
		log.Printf("--- m:\n%v\n\n", string(yaml_temp))
	}
	//TODO Don't Check it

	err = connection.getContents("", "")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(string(connection.ansibleRole.defaults))
	}
}
