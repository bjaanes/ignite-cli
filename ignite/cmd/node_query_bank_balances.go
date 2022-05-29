package ignitecmd

import (
	"fmt"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	sdkflags "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/spf13/cobra"
)

func NewNodeQueryBankBalances() *cobra.Command {
	c := &cobra.Command{
		Use:   "balances [address]",
		Short: "Query for account balances by address",
		RunE:  nodeQueryBankBalancesHandler,
		Args:  cobra.ExactArgs(1),
	}

	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetAccountPrefixes())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	sdkflags.AddPaginationFlagsToCmd(c, "all balances")

	return c
}

func nodeQueryBankBalancesHandler(cmd *cobra.Command, args []string) error {
	var (
		inputAccount   = args[0]
		prefix         = getAddressPrefix(cmd)
		node           = getRPC(cmd)
		home           = getHome(cmd)
		keyringBackend = getKeyringBackend(cmd)
		keyringDir     = getKeyringDir(cmd)
	)
	session := cliui.New()
	defer session.Cleanup()

	client, err := cosmosclient.New(
		cmd.Context(),
		cosmosclient.WithAddressPrefix(prefix),
		cosmosclient.WithHome(home),
		cosmosclient.WithKeyringBackend(keyringBackend),
		cosmosclient.WithKeyringDir(keyringDir),
		cosmosclient.WithNodeAddress(node),
	)
	if err != nil {
		return err
	}

	address, err := client.Bech32Address(inputAccount)
	if err != nil {
		return err
	}

	pagination, err := sdkclient.ReadPageRequest(cmd.Flags())
	if err != nil {
		return err
	}

	session.StartSpinner("Querying")
	balances, err := client.BankBalances(address, pagination)
	if err != nil {
		return err
	}

	var rows [][]string
	for _, b := range balances {
		row := make([]string, 2)
		row[0] = fmt.Sprintf("%s", b.Amount)
		row[1] = b.Denom

		rows = append(rows, row)
	}

	return session.PrintTable([]string{"Amount", "Denom"}, rows...)
}
