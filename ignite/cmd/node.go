package ignitecmd

import (
	"fmt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var rpcAddress string

const (
	flagRpc         = "rpc"
	rpcAddressLocal = "tcp://localhost:26657"

	flagPage       = "page"
	flagLimit      = "limit"
	flagPageKey    = "page-key"
	flagOffset     = "offset"
	flagCountTotal = "count-total"
	flagReverse    = "reverse"
)

func NewNode() *cobra.Command {
	c := &cobra.Command{
		Use:   "node [command]",
		Short: "Make calls to a live blockchain node",
		Args:  cobra.ExactArgs(1),
	}

	c.PersistentFlags().StringVar(&rpcAddress, flagRpc, rpcAddressLocal, "<host>:<port> to tendermint rpc interface for this chain")

	c.AddCommand(NewNodeQuery())
	c.AddCommand(NewNodeTx())

	return c
}

func flagSetPagination(query string) *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.Uint64(flagPage, 1, fmt.Sprintf("pagination page of %s to query for. This sets offset to a multiple of limit", query))
	fs.String(flagPageKey, "", fmt.Sprintf("pagination page-key of %s to query for", query))
	fs.Uint64(flagOffset, 0, fmt.Sprintf("pagination offset of %s to query for", query))
	fs.Uint64(flagLimit, 100, fmt.Sprintf("pagination limit of %s to query for", query))
	fs.Bool(flagCountTotal, false, fmt.Sprintf("count total number of records in %s to query for", query))
	fs.Bool(flagReverse, false, "results are sorted in descending order")

	return fs
}

func getRPC(cmd *cobra.Command) (rpc string) {
	rpc, _ = cmd.Flags().GetString(flagRpc)
	return
}
