package commands

import (
	"fmt"

	"github.com/deislabs/duffle/pkg/action"
	"github.com/deislabs/duffle/pkg/claim"
	"github.com/deislabs/duffle/pkg/credentials"
	"github.com/deislabs/duffle/pkg/utils/crud"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func uninstallCmd(dockerCli command.Cli) *cobra.Command {
	var opts credentialOptions

	cmd := &cobra.Command{
		Use:   "uninstall <installation-name>",
		Short: "Uninstall an application",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall(dockerCli, args[0], opts)
		},
	}
	opts.addFlags(cmd.Flags())

	return cmd
}

func runUninstall(dockerCli command.Cli, claimName string, opts credentialOptions) error {
	defer muteDockerCli(dockerCli)()
	h := duffleHome()

	claimStore := claim.NewClaimStore(crud.NewFileSystemStore(h.Claims(), "json"))
	c, err := claimStore.Read(claimName)
	if err != nil {
		return err
	}
	opts.SetDefaultTargetContext(dockerCli)
	bind, err := requiredClaimBindMount(c, opts.targetContext, dockerCli)
	if err != nil {
		return err
	}
	driverImpl, errBuf, err := prepareDriver(dockerCli, bind, nil)
	if err != nil {
		return err
	}
	creds, err := prepareCredentialSet(c.Bundle, opts.CredentialSetOpts(dockerCli)...)
	if err != nil {
		return err
	}
	if err := credentials.Validate(creds, c.Bundle.Credentials); err != nil {
		return err
	}
	uninst := &action.Uninstall{
		Driver: driverImpl,
	}
	err = uninst.Run(&c, creds, dockerCli.Out())
	if err == nil {
		return claimStore.Delete(claimName)
	}
	if err2 := claimStore.Store(c); err2 != nil {
		fmt.Fprintf(dockerCli.Err(), "failed to update claim: %s\n", err2)
	}
	return fmt.Errorf("uninstall failed: %s", errBuf)
}
