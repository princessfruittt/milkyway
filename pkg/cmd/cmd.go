package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

type cmdsBuilder struct {
	builderCommon
	cmds []cmdErrorHandler
}

func newCommandBuilder() *cmdsBuilder {
	return &cmdsBuilder{}
}

func (b *cmdsBuilder) addCMDS(cmds ...cmdErrorHandler) *cmdsBuilder {
	b.cmds = append(b.cmds, cmds...)
	return b
}
func (b *cmdsBuilder) addAllCMDS(cmds ...cmdErrorHandler) *cmdsBuilder {
	b.addCMDS(b.newGenerateCmd())
	return b
}
func (b *cmdsBuilder) build() *commonCmd {
	c := b.newCommonCMD()
	addCmds(c.getCommand(), b.cmds...)
	return c
}
func addCmds(root *cobra.Command, cmds ...cmdErrorHandler) {
	for _, command := range cmds {
		cmd := command.getCommand()
		if cmd == nil {
			continue
		}
		root.AddCommand(cmd)
	}
}

var _ cmdsBuilderGetter = (*baseBuilderCmd)(nil)

type cmdsBuilderGetter interface {
	getcmdsBuilder() *cmdsBuilder
}
type baseCmd struct {
	cmd *cobra.Command
}

type baseBuilderCmd struct {
	*baseCmd
	*cmdsBuilder
}

func (b baseBuilderCmd) getcmdsBuilder() *cmdsBuilder {
	return b.cmdsBuilder
}
func (c *baseCmd) getCommand() *cobra.Command {
	return c.cmd
}
func newBaseCmd(cmd *cobra.Command) *baseCmd {
	return &baseCmd{cmd: cmd}
}

type commonCmd struct {
	*baseBuilderCmd
	c *cmdError
}

var _ cmdErrorHandler = (*nilCommand)(nil)

func (c *nilCommand) getCommand() *cobra.Command {
	return nil
}

type nilCommand struct {
}

func (b *cmdsBuilder) newCommonCMD() *commonCmd {
	cc := &commonCmd{}
	cc.baseBuilderCmd = b.newBuilderCmd(&cobra.Command{
		Use:                   "milkyway [flags]",
		DisableFlagsInUseLine: true,
		Short:                 "milkyway generates TOSCA node types from Ansible role",
		Long: `"milkyway" is the main command, used to generate TOSCA node type from Ansible Galaxy role.
Use the generated output tosca.node.Type with TOSCA orchestrator"`,
		Example: `milkyway -u "https://github.com/geerlingguy/ansible-role-nginx"`,
		Run: func(cmd *cobra.Command, args []string) {
			defer cc.timeTrack(time.Now(), "Total")
			c := result{"my result"}
			fmt.Print(c.build())
		},
	})

	return cc
}

type builderCommon struct {
	toscaVersion string
	quiet        bool
}

func (cc *builderCommon) handleFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&cc.toscaVersion, "tosca", "t", "2", `Version of TOSCA. e.g: -tv 1. 
List of versions:
	1 - TOSCA v1.0
	1.1 or 1.2 or 1.3 - TOSCA Simple YAML Profile versions
	2 - TOSCA v2.0
`)
}

func (b *cmdsBuilder) newBuilderCmd(cmd *cobra.Command) *baseBuilderCmd {
	bbc := &baseBuilderCmd{cmdsBuilder: b, baseCmd: &baseCmd{cmd: cmd}}
	bbc.builderCommon.handleFlags(cmd)
	return bbc
}
func (cc *builderCommon) timeTrack(start time.Time, name string) {
	if cc.quiet {
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("%s in %v ms\n", name, int(1000*elapsed.Seconds()))
}

//func checkErr(logger loggers.Logger, err error, s ...string) {
//	if err == nil {
//		return
//	}
//	for _, message := range s {
//		logger.Errorln(message)
//	}
//	logger.Errorln(err)
//}
