package ignitecmd

import "github.com/spf13/cobra"

func NewNodeQuery() *cobra.Command {
	c := &cobra.Command{
		Use:     "query",
		Short:   "Querying subcommands",
		Aliases: []string{"q"},
	}

	c.AddCommand(NewNodeQueryBank())

	return c
}
