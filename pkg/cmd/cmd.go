package cmd

import (
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
	//b.addAllCMDS(
	//	b.sjjsjds)
	return b
}
func (b *cmdsBuilder) build() *commonCmd {
	c := b.newCommoncmd()
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

func (b *cmdsBuilder) newBuilderCmd(cmd *cobra.Command) *baseBuilderCmd  {
	bbc := &baseBuilderCmd{cmdsBuilder: b, baseCmd: &baseCmd{cmd: cmd}}
	bbc.builderCommon.handleFlags(cmd)
	return bbc
}

func (b *cmdsBuilder) newBuilderBasicCmd(cmd *cobra.Command) *baseBuilderCmd  {
	bbc := &baseBuilderCmd{cmdsBuilder: b, baseCmd: &baseCmd{cmd: cmd}}
	bbc.builderCommon.handleCommonBuilderFlags(cmd)
	return bbc
}


type commonCmd struct {
	*baseBuilderCmd
}

var _ cmdErrorHandler = (*nilCommand)(nil)

type nilCommand struct {
}
func (c *nilCommand) getCommand() *cobra.Command {
	return nil
}


func (b *cmdsBuilder) newCommonCMD() *commonCmd {
	cc := &commonCmd{}
	cc.baseBuilderCmd = b.newBuilderCmd(&cobra.Command{
		//Scarpia/Talos
		Use:   "amaranth -u URL",
		DisableFlagsInUseLine: true,
		Short: "amaranth generate TOSCA node types from Ansible role",
		Long:  "amaranth is the main command, used to generate TOSCA node types from Ansible Galaxy role",
		Example: `
		# Apply the configuration in pod.json to a pod.
		amaranth -u "https://github.com/gantsign/ansible-role-golang"`,
		Run: func(cmd *cobra.Command, args []string) error {
			defer cc.timeTrack(time.Now(), "Total")
			cc.c =c
			return c.build()
		},
	})

	return cc
}
type builderCommon struct {
	roleURL     string

}
func (cc *commonCmd) handleFlags(cmd *cobra.Command){
	cc.handleFlags(cmd)
	cmd.Flags().StringVarP(&cc.roleURL, "roleURL", "u", "", "Ansible galaxy GitHub URL e.g. https://github.com/gantsign/ansible-role-golang")
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