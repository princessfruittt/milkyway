package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type generateCmd struct {
	*baseBuilderCmd
}

func (b *cmdsBuilder) newGenerateCmd() *generateCmd {
	cc := &generateCmd{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate tosca types",
		Run:   cc.printOutput,
	}
	cc.baseBuilderCmd = b.newBuilderCmd(cmd)
	return cc
}
func (c *generateCmd) printOutput(cmd *cobra.Command, args []string) {
	fmt.Println("I am in generated func")
}

//import (
//	"context"
//	"github.com/google/go-github/v33/github"
//)
//func start() {
//	client := github.NewClient(nil)
//	ctx := context.Background()
//
//	// list all organizations for user "willnorris"
//	opt, _, _ := client.Organizations.List(context.Background(), "willnorris", nil)
//}
