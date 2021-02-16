package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strings"
)

type generateCmd struct {
	*baseBuilderCmd
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
func (c *generateCmd) printOutput(cmd *cobra.Command, args []string) {
	fmt.Println("I am in generated func")
}
func (c *generateCmd) generateTypes(cmd *cobra.Command, args []string) {
	client := github.NewClient(nil)
	ctx := context.Background()
	u, err := url.Parse(c.roleURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.Path)
	info := strings.Split(u.Path, "/")
	opt, _, err := client.Repositories.Get(ctx, info[1], info[2])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(opt.GetSource().Source)
}
