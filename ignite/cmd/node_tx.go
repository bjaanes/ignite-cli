package ignitecmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func NewNodeTx() *cobra.Command {
	c := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	c.AddCommand(NewNodeTxBank())

	return c
}

func flagTxFrom() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.String(flagFrom, "", "Account name to use for sending transactions")
	return fs
}
