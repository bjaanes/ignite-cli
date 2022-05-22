package node_test

import (
	"bytes"
	"context"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/ignite-hq/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"strings"
	"testing"
	"time"

	envtest "github.com/ignite-hq/cli/integration"
	"github.com/stretchr/testify/require"
)

const testPrefix = "testpref"

func TestNodeQueryBankBalances(t *testing.T) {
	var (
		env     = envtest.New(t)
		path    = env.Scaffold("github.com/test/blog", "--address-prefix", testPrefix)
		servers = env.RandomizeServerPorts(path, "")
		homeDir = env.SetRandomHomeConfig(path, "")
	)

	env.SetKeyringBackend(keyring.BackendTest, path, "")

	var (
		ctx, cancel = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	)

	go func() {
		defer cancel()
		isBackendAliveErr := env.IsAppServed(ctx, servers)
		require.NoError(t, isBackendAliveErr, "app cannot get online in time")

		// error "account doesn't have any balances" occurs if a sleep is not included
		time.Sleep(time.Second * 1)

		accounts := env.AccountsInKeyring(homeDir, cosmosaccount.KeyringTest, testPrefix)
		var alice envtest.Account
		for _, acc := range accounts {
			if acc.Name == "alice" {
				alice = acc
				break
			}
		}

		var accountOutputBuffer = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
			)),
			envtest.ExecStdout(accountOutputBuffer),
		))
		require.True(t, strings.Contains(accountOutputBuffer.String(), `Amount 		Denom 	
100000000 	stake 	
20000 		token`))

		var addressOutputBuffer = &bytes.Buffer{}
		env.Must(env.Exec("query bank balances",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					alice.Address,
					"--rpc",
					"http://"+servers.RPC,
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
			)),
			envtest.ExecStdout(addressOutputBuffer),
		))
		require.True(t, strings.Contains(addressOutputBuffer.String(), `Amount 		Denom 	
100000000 	stake 	
20000 		token`))

		env.Must(env.Exec("query bank balances fail with non-existent account name",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"nonexistentaccount",
					"--rpc",
					"http://"+servers.RPC,
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
			)),
			envtest.ExecShouldError(),
		))

		env.Must(env.Exec("query bank balances fail with non-existent address",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					testPrefix+"1gspvt8qsk8cryrsxnqt452cjczjm5ejdgla24e",
					"--rpc",
					"http://"+servers.RPC,
					"--home",
					homeDir,
					"--address-prefix",
					testPrefix,
				),
				step.Workdir(path),
			)),
			envtest.ExecShouldError(),
		))

		env.Must(env.Exec("query bank balances fail with wrong prefix",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"node",
					"query",
					"bank",
					"balances",
					"alice",
					"--rpc",
					"http://"+servers.RPC,
					"--home",
					homeDir,
				),
				step.Workdir(path),
			)),
			envtest.ExecShouldError(),
		))
	}()
	env.Must(env.Serve("should serve with Stargate version", path, "", "", envtest.ExecCtx(ctx)))
}
